package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/Yogdunana/yogduoj/judge/internal/checker"
	"github.com/Yogdunana/yogduoj/judge/internal/compiler"
	"github.com/Yogdunana/yogduoj/judge/internal/runner"
	pb "github.com/Yogdunana/yogduoj/judge/proto"
)

// JudgeTask represents a single judging request to be processed.
type JudgeTask struct {
	Request *pb.JudgeRequest
	Done    chan *pb.JudgeResponse
}

// Worker processes judge tasks from a queue.
type Worker struct {
	id         int
	taskCh     <-chan *JudgeTask
	runner     *runner.Runner
	sandboxDir string
	callbackURL string
	httpClient *http.Client
	running    int32
}

// NewWorker creates a new judge worker.
func NewWorker(id int, taskCh <-chan *JudgeTask, sandboxDir string, callbackURL string) *Worker {
	return &Worker{
		id:          id,
		taskCh:      taskCh,
		runner:      runner.NewRunner(sandboxDir),
		sandboxDir:  sandboxDir,
		callbackURL: callbackURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Start begins the worker's main loop, processing tasks from the channel.
func (w *Worker) Start() {
	log.Printf("Worker %d started", w.id)
	go w.run()
}

// IsRunning returns whether the worker is currently processing a task.
func (w *Worker) IsRunning() bool {
	return atomic.LoadInt32(&w.running) == 1
}

// run is the main worker loop.
func (w *Worker) run() {
	for task := range w.taskCh {
		atomic.StoreInt32(&w.running, 1)
		response := w.processTask(task.Request)
		atomic.StoreInt32(&w.running, 0)

		// Send the response back through the done channel.
		if task.Done != nil {
			task.Done <- response
		}

		// Also send callback to backend.
		w.sendCallback(response)
	}
	log.Printf("Worker %d stopped", w.id)
}

// processTask handles a single judge request end-to-end.
func (w *Worker) processTask(req *pb.JudgeRequest) *pb.JudgeResponse {
	log.Printf("Worker %d: processing submission %s (language=%s, judge_type=%s)",
		w.id, req.SubmissionId, req.Language, req.JudgeType)

	// Create a temporary working directory for this submission.
	workDir, err := os.MkdirTemp(w.sandboxDir, fmt.Sprintf("sub-%s-", req.SubmissionId))
	if err != nil {
		return w.makeErrorResponse(req.SubmissionId, "SE", 0, 0,
			fmt.Sprintf("failed to create work directory: %v", err))
	}
	defer os.RemoveAll(workDir)

	// Write source code to file.
	sourcePath := filepath.Join(workDir, sourceFilename(req.Language))
	if err := os.WriteFile(sourcePath, []byte(req.Code), 0644); err != nil {
		return w.makeErrorResponse(req.SubmissionId, "SE", 0, 0,
			fmt.Sprintf("failed to write source file: %v", err))
	}

	// Step 1: Compile the code.
	comp := compiler.GetCompiler(req.Language)
	if comp == nil {
		return w.makeErrorResponse(req.SubmissionId, "CE", 0, 0,
			fmt.Sprintf("unsupported language: %s", req.Language))
	}

	compileTimeLimit := time.Duration(req.TimeLimitMs*2) * time.Millisecond
	compileMemLimit := req.MemoryLimitKb * 2
	compileResult := comp.Compile(sourcePath, workDir, compileTimeLimit, compileMemLimit)

	if !compileResult.Success {
		return w.makeErrorResponse(req.SubmissionId, "CE", 0, 0, compileResult.Error)
	}

	// Step 2: Determine the checker.
	var chk checker.Checker
	switch req.JudgeType {
	case "ctf":
		chk = &checker.CTFChecker{}
	default:
		chk = &checker.StandardChecker{}
	}

	// Step 3: Run each test case.
	caseResults := make([]*pb.TestCaseResult, len(req.TestCases))
	totalScore := int32(0)
	maxTime := int64(0)
	maxMemory := int64(0)
	allAccepted := true

	for i, tc := range req.TestCases {
		caseResult := w.runTestCase(
			req.SubmissionId,
			req.Language,
			compileResult.OutputPath,
			tc,
			i,
			req.TimeLimitMs,
			req.MemoryLimitKb,
			chk,
			req.CtfFlag,
			req.JudgeType,
		)
		caseResults[i] = caseResult

		// Track maximum resource usage.
		if caseResult.TimeUsedMs > maxTime {
			maxTime = caseResult.TimeUsedMs
		}
		if caseResult.MemoryUsedKb > maxMemory {
			maxMemory = caseResult.MemoryUsedKb
		}

		// Accumulate score.
		totalScore += caseResult.Score

		// Check if this case failed.
		if caseResult.Verdict != string(checker.VerdictAC) {
			allAccepted = false
		}
	}

	// Determine overall verdict.
	overallVerdict := "AC"
	if !allAccepted {
		overallVerdict = "WA"
		// Check for specific error verdicts.
		for _, cr := range caseResults {
			switch cr.Verdict {
			case "CE", "RE", "TLE", "MLE":
				overallVerdict = cr.Verdict
				break
			}
		}
	}

	// For OI mode, total_score is the sum of individual case scores.
	// For ACM mode, total_score is either 100 (AC) or 0.
	if req.JudgeType != "oi" {
		if allAccepted {
			totalScore = 100
		} else {
			totalScore = 0
		}
	}

	log.Printf("Worker %d: submission %s completed with verdict=%s, score=%d",
		w.id, req.SubmissionId, overallVerdict, totalScore)

	return &pb.JudgeResponse{
		SubmissionId: req.SubmissionId,
		Verdict:      overallVerdict,
		TotalScore:   totalScore,
		TimeUsedMs:   maxTime,
		MemoryUsedKb: maxMemory,
		CaseResults:  caseResults,
	}
}

// runTestCase runs a single test case and returns the result.
func (w *Worker) runTestCase(
	submissionID string,
	language string,
	execPath string,
	tc *pb.TestCase,
	index int,
	timeLimitMs int64,
	memoryLimitKB int64,
	chk checker.Checker,
	ctfFlag string,
	judgeType string,
) *pb.TestCaseResult {
	containerID := fmt.Sprintf("judge-run-%s-%d", submissionID, index)

	// Determine input for the test case.
	input := tc.Input
	if judgeType == "ctf" {
		// For CTF, input is not used; we just check the output against the flag.
		input = ""
	}

	// Run the program.
	result := w.runner.RunWithCgroup(
		containerID,
		execPath,
		language,
		input,
		timeLimitMs,
		memoryLimitKB,
		1.0, // 1 CPU core
		50,  // max 50 processes
	)

	caseResult := &pb.TestCaseResult{
		Index:        int32(index),
		TimeUsedMs:   result.TimeUsedMs,
		MemoryUsedKb: result.MemoryUsedKb,
	}

	// Determine verdict based on execution result.
	if result.TimedOut {
		caseResult.Verdict = "TLE"
		caseResult.ErrorMessage = "Time Limit Exceeded"
		caseResult.Score = 0
		return caseResult
	}

	if result.OOMKilled {
		caseResult.Verdict = "MLE"
		caseResult.ErrorMessage = "Memory Limit Exceeded"
		caseResult.Score = 0
		return caseResult
	}

	if result.ExitCode != 0 {
		caseResult.Verdict = "RE"
		caseResult.ErrorMessage = fmt.Sprintf("Runtime Error (exit code %d): %s", result.ExitCode, result.Stderr)
		caseResult.Score = 0
		return caseResult
	}

	// Check the output.
	var verdict checker.Verdict
	if judgeType == "ctf" {
		// For CTF, compare stdout against the expected flag.
		ctfChecker := &checker.CTFChecker{}
		verdict = ctfChecker.Check(result.Stdout, ctfFlag)
	} else {
		verdict = chk.Check(result.Stdout, tc.ExpectedOutput)
	}

	caseResult.Verdict = string(verdict)

	// Calculate score based on weight.
	if verdict == checker.VerdictAC || verdict == checker.VerdictPE {
		caseResult.Score = tc.ScoreWeight
	} else {
		caseResult.Score = 0
	}

	return caseResult
}

// makeErrorResponse creates a JudgeResponse for system/compilation errors.
func (w *Worker) makeErrorResponse(submissionID string, verdict string, timeMs int64, memKb int64, errMsg string) *pb.JudgeResponse {
	return &pb.JudgeResponse{
		SubmissionId: submissionID,
		Verdict:      verdict,
		TotalScore:   0,
		TimeUsedMs:   timeMs,
		MemoryUsedKb: memKb,
		ErrorMessage: errMsg,
		CaseResults:  nil,
	}
}

// sendCallback sends the judge result back to the backend via HTTP POST.
func (w *Worker) sendCallback(response *pb.JudgeResponse) {
	if w.callbackURL == "" {
		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("Worker %d: failed to marshal callback response: %v", w.id, err)
		return
	}

	resp, err := w.httpClient.Post(w.callbackURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("Worker %d: failed to send callback to %s: %v", w.id, w.callbackURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Worker %d: callback returned status %d for submission %s",
			w.id, resp.StatusCode, response.SubmissionId)
	}
}

// sourceFilename returns the appropriate source filename for a language.
func sourceFilename(language string) string {
	switch language {
	case "cpp", "c++", "cc", "cxx":
		return "solution.cpp"
	case "c":
		return "solution.c"
	case "java":
		return "Main.java"
	case "python", "python3", "py":
		return "solution.py"
	default:
		return "solution"
	}
}
