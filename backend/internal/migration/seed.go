package migration

import (
	"fmt"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedData populates the database with initial seed data for development/demo.
func SeedData(db *gorm.DB) error {
	// --- Tags ---
	tags := []model.Tag{
		{Name: "数组"}, {Name: "字符串"}, {Name: "动态规划"}, {Name: "贪心"},
		{Name: "图论"}, {Name: "二分查找"}, {Name: "排序"}, {Name: "搜索"},
		{Name: "数学"}, {Name: "模拟"}, {Name: "栈"}, {Name: "队列"},
		{Name: "树"}, {Name: "哈希表"}, {Name: "DFS"}, {Name: "BFS"},
		{Name: "数论"}, {Name: "位运算"},
	}
	for i := range tags {
		db.FirstOrCreate(&tags[i], model.Tag{Name: tags[i].Name})
	}

	// --- Users ---
	adminPassHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userPassHash, _ := bcrypt.GenerateFromPassword([]byte("user123456"), bcrypt.DefaultCost)

	now := time.Now()
	users := []model.User{
		{
			Username: "admin", Email: "admin@yogdunana.com", PasswordHash: string(adminPassHash),
			Role: "admin", Status: "active", SolvedCount: 42, SubmissionCount: 87, ContestCount: 5,
			Rating: 1850, CreatedAt: now.AddDate(-1, 0, 0),
		},
		{
			Username: "zhangsan", Email: "zhangsan@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 35, SubmissionCount: 72, ContestCount: 4,
			Rating: 1720, CreatedAt: now.AddDate(-10, 0, 0),
		},
		{
			Username: "lisi", Email: "lisi@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 28, SubmissionCount: 55, ContestCount: 3,
			Rating: 1650, CreatedAt: now.AddDate(-9, 0, 0),
		},
		{
			Username: "wangwu", Email: "wangwu@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 51, SubmissionCount: 103, ContestCount: 6,
			Rating: 1980, CreatedAt: now.AddDate(-8, 0, 0),
		},
		{
			Username: "chenliu", Email: "chenliu@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 19, SubmissionCount: 41, ContestCount: 2,
			Rating: 1480, CreatedAt: now.AddDate(-7, 0, 0),
		},
		{
			Username: "sunqi", Email: "sunqi@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 63, SubmissionCount: 120, ContestCount: 8,
			Rating: 2100, CreatedAt: now.AddDate(-6, 0, 0),
		},
		{
			Username: "zhouba", Email: "zhouba@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 12, SubmissionCount: 28, ContestCount: 1,
			Rating: 1320, CreatedAt: now.AddDate(-5, 0, 0),
		},
		{
			Username: "alice", Email: "alice@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 40, SubmissionCount: 85, ContestCount: 5,
			Rating: 1780, CreatedAt: now.AddDate(-4, 0, 0),
		},
		{
			Username: "bob", Email: "bob@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 8, SubmissionCount: 15, ContestCount: 1,
			Rating: 1200, CreatedAt: now.AddDate(-3, 0, 0),
		},
		{
			Username: "charlie", Email: "charlie@example.com", PasswordHash: string(userPassHash),
			Role: "user", Status: "active", SolvedCount: 22, SubmissionCount: 48, ContestCount: 3,
			Rating: 1550, CreatedAt: now.AddDate(-2, 0, 0),
		},
	}
	for i := range users {
		db.FirstOrCreate(&users[i], model.User{Username: users[i].Username})
	}

	// Helper: get tag by name
	tagMap := make(map[string]model.Tag)
	for _, t := range tags {
		tagMap[t.Name] = t
	}

	// --- Problems ---
	problems := []model.Problem{
		{
			Title: "A+B Problem", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定两个整数 $A$ 和 $B$，请计算 $A + B$ 的值。\n\n## 输入格式\n\n一行，包含两个整数 $A$ 和 $B$，用空格分隔。\n\n## 输出格式\n\n一行，输出 $A + B$ 的结果。\n\n## 样例\n\n```plain\n输入：\n1 2\n\n输出：\n3\n```\n\n## 数据范围\n\n$-10^9 \\le A, B \\le 10^9$",
			InputFormat:  "一行两个整数 A B",
			OutputFormat: "一行一个整数，表示 A+B 的结果",
			Hints:        "直接使用加法运算即可。",
			Source:       "original",
			SubmitCount:  120, AcceptedCount: 95, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-10, 0, 0),
		},
		{
			Title: "两数之和", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个整数数组 `nums` 和一个整数目标值 `target`，请你在该数组中找出和为目标值 `target` 的那两个整数，并返回它们的数组下标。\n\n你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。\n\n## 输入格式\n\n第一行一个整数 $n$，表示数组长度。\n第二行 $n$ 个整数，表示数组元素。\n第三行一个整数 $target$，表示目标和。\n\n## 输出格式\n\n一行两个整数，表示两个元素的下标（从 0 开始），按升序输出。\n\n## 样例\n\n```plain\n输入：\n4\n2 7 11 15\n9\n\n输出：\n0 1\n```\n\n## 数据范围\n\n$2 \\le n \\le 10^4$\n$-10^9 \\le nums[i] \\le 10^9$\n$-10^9 \\le target \\le 10^9$",
			InputFormat:  "第一行 n，第二行 n 个整数，第三行 target",
			OutputFormat: "两个下标，空格分隔",
			Hints:        "考虑使用哈希表来优化查找过程。",
			Source:       "original",
			SubmitCount:  85, AcceptedCount: 62, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-9, -6, 0),
		},
		{
			Title: "最长回文子串", Type: "programming", Difficulty: "medium",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个字符串 `s`，请你找出其中最长的回文子串。\n\n如果字符串的反序与原始字符串相同，则该字符串称为回文字符串。\n\n## 输入格式\n\n一行一个字符串 $s$，仅包含小写字母和数字。\n\n## 输出格式\n\n一行一个字符串，表示最长的回文子串。如果有多个，输出字典序最小的那个。\n\n## 样例\n\n```plain\n输入：\nbabad\n\n输出：\nbab\n```\n\n```plain\n输入：\ncbbd\n\n输出：\nbb\n```\n\n## 数据范围\n\n$1 \\le |s| \\le 1000$",
			InputFormat:  "一行一个字符串 s",
			OutputFormat: "最长回文子串",
			Hints:        "可以使用动态规划或中心扩展法。",
			Source:       "original",
			SubmitCount:  67, AcceptedCount: 38, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-9, -3, 0),
		},
		{
			Title: "合并两个有序数组", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给你两个按非递减顺序排列的整数数组 `nums1` 和 `nums2`，另有两个整数 `m` 和 `n`，分别表示 `nums1` 和 `nums2` 中的元素数目。\n\n请你合并 `nums2` 到 `nums1` 中，使合并后的数组同样按非递减顺序排列。\n\n## 输入格式\n\n第一行三个整数 $m, n$，分别表示两个数组的长度。\n第二行 $m$ 个非递减整数，表示 `nums1`。\n第三行 $n$ 个非递减整数，表示 `nums2`。\n\n## 输出格式\n\n一行 $m+n$ 个非递减整数，表示合并后的数组。\n\n## 样例\n\n```plain\n输入：\n3 3\n1 2 3\n2 5 6\n\n输出：\n1 2 2 3 5 6\n```\n\n## 数据范围\n\n$0 \\le m, n \\le 200$\n$-10^9 \\le nums[i] \\le 10^9$",
			InputFormat:  "第一行 m n，第二行 nums1，第三行 nums2",
			OutputFormat: "合并后的有序数组",
			Hints:        "从后向前合并可以避免额外空间。",
			Source:       "original",
			SubmitCount:  55, AcceptedCount: 45, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-8, -6, 0),
		},
		{
			Title: "二叉树的最大深度", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个二叉树，找出其最大深度。\n\n二叉树的深度为根节点到最远叶子节点的最长路径上的节点数。\n\n## 输入格式\n\n第一行一个整数 $n$，表示节点数。\n接下来 $n$ 行，每行两个整数 $left$ 和 $right$，表示第 $i$ 个节点的左右子节点编号（0 表示空节点）。节点编号从 1 开始。\n\n## 输出格式\n\n一个整数，表示二叉树的最大深度。\n\n## 样例\n\n```plain\n输入：\n3\n2 3\n0 0\n0 0\n\n输出：\n2\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^4$",
			InputFormat:  "第一行 n，接下来 n 行每行两个整数表示左右子节点",
			OutputFormat: "最大深度",
			Hints:        "使用递归或 BFS 层序遍历。",
			Source:       "original",
			SubmitCount:  48, AcceptedCount: 40, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-8, -3, 0),
		},
		{
			Title: "爬楼梯", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n假设你正在爬楼梯。需要 $n$ 阶你才能到达楼顶。\n\n每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？\n\n## 输入格式\n\n一个整数 $n$。\n\n## 输出格式\n\n一个整数，表示不同的方法数。由于结果可能很大，请对 $10^9 + 7$ 取模。\n\n## 样例\n\n```plain\n输入：\n3\n\n输出：\n3\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^5$",
			InputFormat:  "一个整数 n",
			OutputFormat: "方法数（对 1e9+7 取模）",
			Hints:        "这就是斐波那契数列的变体。",
			Source:       "original",
			SubmitCount:  72, AcceptedCount: 60, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-7, -6, 0),
		},
		{
			Title: "最长递增子序列", Type: "programming", Difficulty: "medium",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给你一个整数数组 `nums`，找到其中最长严格递增子序列的长度。\n\n子序列是由数组派生而来的序列，删除（或不删除）数组中的元素而不改变其余元素的顺序。\n\n## 输入格式\n\n第一行一个整数 $n$。\n第二行 $n$ 个整数。\n\n## 输出格式\n\n一个整数，表示最长递增子序列的长度。\n\n## 样例\n\n```plain\n输入：\n10\n10 9 2 5 3 7 101 18\n\n输出：\n4\n```\n\n## 数据范围\n\n$1 \\le n \\le 2500$\n$-10^4 \\le nums[i] \\le 10^4$",
			InputFormat:  "第一行 n，第二行 n 个整数",
			OutputFormat: "最长递增子序列的长度",
			Hints:        "可以使用动态规划 O(n^2) 或贪心+二分 O(n log n)。",
			Source:       "original",
			SubmitCount:  58, AcceptedCount: 30, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-7, -3, 0),
		},
		{
			Title: "最小生成树", Type: "programming", Difficulty: "medium",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个无向连通图，请你求出它的最小生成树的边权之和。\n\n## 输入格式\n\n第一行两个整数 $n, m$，分别表示点数和边数。\n接下来 $m$ 行，每行三个整数 $u, v, w$，表示 $u$ 和 $v$ 之间有一条权值为 $w$ 的边。\n\n## 输出格式\n\n一个整数，表示最小生成树的边权之和。如果图不连通，输出 `-1`。\n\n## 样例\n\n```plain\n输入：\n4 5\n1 2 1\n1 3 2\n2 3 3\n2 4 4\n3 4 5\n\n输出：\n7\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^5$\n$1 \\le m \\le 2 \\times 10^5$\n$1 \\le w \\le 10^4$",
			InputFormat:  "第一行 n m，接下来 m 行每行 u v w",
			OutputFormat: "最小生成树的边权之和",
			Hints:        "使用 Kruskal 算法配合并查集，或 Prim 算法。",
			Source:       "original",
			SubmitCount:  42, AcceptedCount: 18, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-6, -6, 0),
		},
		{
			Title: "单源最短路", Type: "programming", Difficulty: "medium",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个 $n$ 个点 $m$ 条边的有向图，求从点 1 到点 $n$ 的最短路径长度。\n\n## 输入格式\n\n第一行两个整数 $n, m$。\n接下来 $m$ 行，每行三个整数 $u, v, w$，表示从 $u$ 到 $v$ 有一条权值为 $w$ 的有向边。\n\n## 输出格式\n\n一个整数，表示从 1 到 $n$ 的最短路径长度。如果不存在路径，输出 `-1`。\n\n## 样例\n\n```plain\n输入：\n3 3\n1 2 1\n2 3 2\n1 3 5\n\n输出：\n3\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^5$\n$1 \\le m \\le 2 \\times 10^5$\n$1 \\le w \\le 10^4$",
			InputFormat:  "第一行 n m，接下来 m 行每行 u v w",
			OutputFormat: "最短路径长度，不存在则输出 -1",
			Hints:        "使用 Dijkstra 算法，注意边权非负。",
			Source:       "original",
			SubmitCount:  50, AcceptedCount: 25, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-6, -3, 0),
		},
		{
			Title: "背包问题", Type: "programming", Difficulty: "medium",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n有 $n$ 件物品和一个容量为 $W$ 的背包。第 $i$ 件物品的重量为 $w_i$，价值为 $v_i$。\n\n求解将哪些物品装入背包，可使这些物品的总重量不超过背包容量，且总价值最大。\n\n## 输入格式\n\n第一行两个整数 $n, W$。\n接下来 $n$ 行，每行两个整数 $w_i, v_i$。\n\n## 输出格式\n\n一个整数，表示最大总价值。\n\n## 样例\n\n```plain\n输入：\n4 5\n1 2\n2 4\n3 4\n4 5\n\n输出：\n8\n```\n\n## 数据范围\n\n$1 \\le n \\le 100$\n$1 \\le W \\le 1000$\n$1 \\le w_i, v_i \\le 1000$",
			InputFormat:  "第一行 n W，接下来 n 行每行 w_i v_i",
			OutputFormat: "最大总价值",
			Hints:        "经典 0-1 背包问题，使用动态规划。",
			Source:       "original",
			SubmitCount:  65, AcceptedCount: 35, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-5, -6, 0),
		},
		{
			Title: "快速排序", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定 $n$ 个整数，请使用快速排序算法将它们按非递减顺序排列。\n\n## 输入格式\n\n第一行一个整数 $n$。\n第二行 $n$ 个整数。\n\n## 输出格式\n\n一行 $n$ 个非递减整数。\n\n## 样例\n\n```plain\n输入：\n5\n3 1 4 1 5\n\n输出：\n1 1 3 4 5\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^5$\n$-10^9 \\le a_i \\le 10^9$",
			InputFormat:  "第一行 n，第二行 n 个整数",
			OutputFormat: "排序后的 n 个整数",
			Hints:        "注意选择合适的 pivot 避免退化。",
			Source:       "original",
			SubmitCount:  90, AcceptedCount: 78, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-5, -3, 0),
		},
		{
			Title: "矩阵快速幂", Type: "programming", Difficulty: "hard",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定 $n \\times n$ 矩阵 $A$ 和非负整数 $k$，求 $A^k$ 对 $10^9+7$ 取模的结果。\n\n## 输入格式\n\n第一行两个整数 $n, k$。\n接下来 $n$ 行，每行 $n$ 个整数，表示矩阵 $A$。\n\n## 输出格式\n\n$n$ 行，每行 $n$ 个整数，表示 $A^k \\bmod (10^9+7)$ 的结果。\n\n## 样例\n\n```plain\n输入：\n2 2\n1 1\n1 0\n\n输出：\n2 1\n1 1\n```\n\n## 数据范围\n\n$1 \\le n \\le 100$\n$0 \\le k \\le 10^{18}$\n$0 \\le A_{ij} \\le 10^9$",
			InputFormat:  "第一行 n k，接下来 n 行 n 列矩阵",
			OutputFormat: "A^k mod 1e9+7 的矩阵",
			Hints:        "使用快速幂的思想，将矩阵乘法和二进制拆分结合。",
			Source:       "original",
			SubmitCount:  30, AcceptedCount: 10, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-4, -6, 0),
		},
		{
			Title: "线段树区间求和", Type: "programming", Difficulty: "hard",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个长度为 $n$ 的数组，支持两种操作：\n1. 将第 $i$ 个元素加上 $d$\n2. 查询区间 $[l, r]$ 的和\n\n## 输入格式\n\n第一行两个整数 $n, q$，表示数组长度和操作数。\n第二行 $n$ 个整数，表示初始数组。\n接下来 $q$ 行，每行一个操作：\n- `1 i d`：将第 $i$ 个元素加上 $d$\n- `2 l r`：查询区间 $[l, r]$ 的和\n\n## 输出格式\n\n对于每个查询操作，输出一行一个整数。\n\n## 样例\n\n```plain\n输入：\n5 3\n1 2 3 4 5\n2 1 3\n1 2 3\n2 1 3\n\n输出：\n6\n9\n```\n\n## 数据范围\n\n$1 \\le n, q \\le 10^5$\n$-10^9 \\le a_i, d \\le 10^9$",
			InputFormat:  "第一行 n q，第二行数组，接下来 q 行操作",
			OutputFormat: "每个查询输出区间和",
			Hints:        "使用线段树或树状数组（BIT）。",
			Source:       "original",
			SubmitCount:  25, AcceptedCount: 8, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-4, -3, 0),
		},
		{
			Title: "字符串匹配 KMP", Type: "programming", Difficulty: "medium",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定两个字符串 `text` 和 `pattern`，请使用 KMP 算法找出 `pattern` 在 `text` 中所有出现的位置。\n\n## 输入格式\n\n第一行一个字符串 `text`。\n第二行一个字符串 `pattern`。\n\n## 输出格式\n\n第一行一个整数 $k$，表示出现次数。\n第二行 $k$ 个整数，表示所有匹配的起始位置（从 0 开始）。\n\n## 样例\n\n```plain\n输入：\nabababab\nabab\n\n输出：\n2\n0 2\n```\n\n## 数据范围\n\n$1 \\le |text| \\le 10^6$\n$1 \\le |pattern| \\le 10^6$",
			InputFormat:  "第一行 text，第二行 pattern",
			OutputFormat: "第一行出现次数，第二行所有起始位置",
			Hints:        "先构造 next 数组（前缀函数），再进行匹配。",
			Source:       "original",
			SubmitCount:  35, AcceptedCount: 15, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-3, -6, 0),
		},
		{
			Title: "N 皇后问题", Type: "programming", Difficulty: "hard",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n在 $n \\times n$ 的棋盘上放置 $n$ 个皇后，使得它们互不攻击（即任意两个皇后不在同一行、同一列、同一对角线上）。\n\n请输出所有不同的解法数量。\n\n## 输入格式\n\n一个整数 $n$。\n\n## 输出格式\n\n一个整数，表示解法数量。\n\n## 样例\n\n```plain\n输入：\n4\n\n输出：\n2\n```\n\n## 数据范围\n\n$1 \\le n \\le 14$",
			InputFormat:  "一个整数 n",
			OutputFormat: "解法数量",
			Hints:        "使用回溯法，注意剪枝优化。",
			Source:       "original",
			SubmitCount:  40, AcceptedCount: 12, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-3, -3, 0),
		},
		{
			Title: "最大子数组和", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个整数数组 `nums`，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。\n\n## 输入格式\n\n第一行一个整数 $n$。\n第二行 $n$ 个整数。\n\n## 输出格式\n\n一个整数，表示最大子数组和。\n\n## 样例\n\n```plain\n输入：\n9\n-2 1 -3 4 -1 2 1 -5 4\n\n输出：\n6\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^5$\n$-10^4 \\le nums[i] \\le 10^4$",
			InputFormat:  "第一行 n，第二行 n 个整数",
			OutputFormat: "最大子数组和",
			Hints:        "使用 Kadane 算法，维护当前子数组和。",
			Source:       "original",
			SubmitCount:  80, AcceptedCount: 68, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-2, -6, 0),
		},
		{
			Title: "二分查找", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个已排序的整数数组和目标值，如果目标值存在于数组中，返回其下标；否则返回 `-1`。\n\n## 输入格式\n\n第一行一个整数 $n$。\n第二行 $n$ 个非递减整数。\n第三行一个整数 $target$。\n\n## 输出格式\n\n一个整数，表示目标值的下标。如果不存在，输出 `-1`。\n\n## 样例\n\n```plain\n输入：\n5\n1 3 5 7 9\n5\n\n输出：\n2\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^5$\n$-10^9 \\le nums[i], target \\le 10^9$",
			InputFormat:  "第一行 n，第二行数组，第三行 target",
			OutputFormat: "目标值下标，不存在输出 -1",
			Hints:        "使用二分查找，注意边界条件。",
			Source:       "original",
			SubmitCount:  95, AcceptedCount: 82, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-2, -3, 0),
		},
		{
			Title: "图的着色", Type: "programming", Difficulty: "hard",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个无向图，判断是否可以用 $k$ 种颜色对图进行着色，使得相邻节点颜色不同。\n\n## 输入格式\n\n第一行三个整数 $n, m, k$，分别表示点数、边数和颜色数。\n接下来 $m$ 行，每行两个整数 $u, v$，表示一条边。\n\n## 输出格式\n\n如果可以着色，输出 `YES`；否则输出 `NO`。\n\n## 样例\n\n```plain\n输入：\n3 3 3\n1 2\n2 3\n1 3\n\n输出：\nYES\n```\n\n## 数据范围\n\n$1 \\le n \\le 20$\n$0 \\le m \\le n(n-1)/2$\n$1 \\le k \\le n$",
			InputFormat:  "第一行 n m k，接下来 m 行边",
			OutputFormat: "YES 或 NO",
			Hints:        "使用回溯法尝试每种颜色。",
			Source:       "original",
			SubmitCount:  20, AcceptedCount: 6, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-1, -6, 0),
		},
		{
			Title: "滑动窗口最大值", Type: "programming", Difficulty: "hard",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给你一个整数数组 `nums`，有一个大小为 $k$ 的滑动窗口从数组的最左侧移动到数组的最右侧。你只可以看到在滑动窗口内的 $k$ 个数字。\n\n请输出每次滑动窗口中的最大值。\n\n## 输入格式\n\n第一行两个整数 $n, k$。\n第二行 $n$ 个整数。\n\n## 输出格式\n\n一行 $n-k+1$ 个整数，表示每个滑动窗口的最大值。\n\n## 样例\n\n```plain\n输入：\n8 3\n1 3 -1 -3 5 3 6 7\n\n输出：\n3 3 5 5 6 7\n```\n\n## 数据范围\n\n$1 \\le k \\le n \\le 10^5$\n$-10^4 \\le nums[i] \\le 10^4$",
			InputFormat:  "第一行 n k，第二行 n 个整数",
			OutputFormat: "每个滑动窗口的最大值",
			Hints:        "使用单调队列维护窗口内的最大值。",
			Source:       "original",
			SubmitCount:  28, AcceptedCount: 9, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(-1, -3, 0),
		},
		{
			Title: "哈希表实现", Type: "programming", Difficulty: "easy",
			TimeLimitMs: 1000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n实现一个简单的哈希表，支持以下操作：\n1. `insert key value`：插入键值对\n2. `get key`：查询键对应的值\n3. `delete key`：删除键\n\n## 输入格式\n\n第一行一个整数 $q$，表示操作数。\n接下来 $q$ 行，每行一个操作。\n\n## 输出格式\n\n对于每个 `get` 操作，输出对应的值。如果键不存在，输出 `null`。\n\n## 样例\n\n```plain\n输入：\n5\ninsert 1 10\ninsert 2 20\nget 1\nget 3\ndelete 1\nget 1\n\n输出：\n10\nnull\nnull\n```\n\n## 数据范围\n\n$1 \\le q \\le 10^5$\n键和值均为整数，$-10^9 \\le key, value \\le 10^9$",
			InputFormat:  "第一行 q，接下来 q 行操作",
			OutputFormat: "每个 get 操作输出对应值或 null",
			Hints:        "使用开放寻址法或链地址法处理冲突。",
			Source:       "original",
			SubmitCount:  45, AcceptedCount: 32, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(0, -1, 0),
		},
		{
			Title: "拓扑排序", Type: "programming", Difficulty: "medium",
			TimeLimitMs: 2000, MemoryLimitKb: 262144,
			Description: "## 题目描述\n\n给定一个有向图，请判断图中是否存在环，如果不存在，输出任意一种拓扑排序结果。\n\n## 输入格式\n\n第一行两个整数 $n, m$，分别表示点数和边数。\n接下来 $m$ 行，每行两个整数 $u, v$，表示一条从 $u$ 到 $v$ 的有向边。\n\n## 输出格式\n\n如果存在环，输出 `-1`。\n否则输出一行 $n$ 个整数，表示拓扑排序结果。\n\n## 样例\n\n```plain\n输入：\n4 4\n1 2\n1 3\n2 4\n3 4\n\n输出：\n1 2 3 4\n```\n\n## 数据范围\n\n$1 \\le n \\le 10^5$\n$0 \\le m \\le 2 \\times 10^5$",
			InputFormat:  "第一行 n m，接下来 m 行边",
			OutputFormat: "拓扑排序结果或 -1",
			Hints:        "使用 Kahn 算法（BFS）或 DFS 检测环。",
			Source:       "original",
			SubmitCount:  38, AcceptedCount: 20, Status: "public", CreatedBy: 1,
			CreatedAt: now.AddDate(0, -1, 0),
		},
	}

	for i := range problems {
		db.FirstOrCreate(&problems[i], model.Problem{Title: problems[i].Title})
	}

	// --- Problem-Tag associations ---
	problemTags := []struct {
		problemIdx int
		tagNames   []string
	}{
		{0, []string{"数学", "模拟"}},
		{1, []string{"数组", "哈希表"}},
		{2, []string{"字符串", "动态规划"}},
		{3, []string{"数组", "排序"}},
		{4, []string{"树", "DFS", "BFS"}},
		{5, []string{"动态规划", "数学"}},
		{6, []string{"动态规划", "二分查找"}},
		{7, []string{"图论", "贪心"}},
		{8, []string{"图论", "贪心"}},
		{9, []string{"动态规划"}},
		{10, []string{"排序"}},
		{11, []string{"数学", "动态规划"}},
		{12, []string{"树", "动态规划"}},
		{13, []string{"字符串"}},
		{14, []string{"搜索", "DFS"}},
		{15, []string{"数组", "动态规划"}},
		{16, []string{"二分查找", "数组"}},
		{17, []string{"图论", "搜索", "DFS"}},
		{18, []string{"队列", "数组"}},
		{19, []string{"哈希表"}},
		{20, []string{"图论", "搜索", "BFS"}},
	}

	for _, pt := range problemTags {
		if pt.problemIdx >= len(problems) {
			continue
		}
		var existingTags []model.Tag
		db.Model(&problems[pt.problemIdx]).Association("Tags").Find(&existingTags)
		existingTagMap := make(map[string]bool)
		for _, t := range existingTags {
			existingTagMap[t.Name] = true
		}
		for _, tn := range pt.tagNames {
			if !existingTagMap[tn] {
				if t, ok := tagMap[tn]; ok {
					db.Model(&problems[pt.problemIdx]).Association("Tags").Append(&t)
				}
			}
		}
	}

	// --- Samples for first 5 problems ---
	sampleData := []struct {
		problemIdx int
		input      string
		output     string
		order      int
	}{
		{0, "1 2\n", "3\n", 1},
		{0, "100 200\n", "300\n", 2},
		{1, "4\n2 7 11 15\n9\n", "0 1\n", 1},
		{1, "3\n3 2 4\n6\n", "0 2\n", 2},
		{2, "babad\n", "bab\n", 1},
		{2, "cbbd\n", "bb\n", 2},
		{3, "3 3\n1 2 3\n2 5 6\n", "1 2 2 3 5 6\n", 1},
		{4, "3\n2 3\n0 0\n0 0\n", "2\n", 1},
	}

	for _, sd := range sampleData {
		if sd.problemIdx >= len(problems) {
			continue
		}
		pid := problems[sd.problemIdx].ID

		// Create TestData first
		td := model.TestData{
			ProblemID:  pid,
			InputFile:  fmt.Sprintf("sample_%d.in", sd.order),
			OutputFile: fmt.Sprintf("sample_%d.out", sd.order),
			IsSample:   true,
			Generation: "manual",
		}
		db.FirstOrCreate(&td, model.TestData{
			ProblemID: pid,
			InputFile: td.InputFile,
		})

		// Create Sample linking
		sample := model.Sample{
			ProblemID:    pid,
			TestDataID:   td.ID,
			DisplayOrder: sd.order,
		}
		db.FirstOrCreate(&sample, model.Sample{
			ProblemID:    pid,
			TestDataID:   td.ID,
			DisplayOrder: sd.order,
		})
	}

	// --- Submissions ---
	acCodeSnippets := map[string]string{
		"cpp": `#include <iostream>
using namespace std;
int main() {
    int a, b;
    cin >> a >> b;
    cout << a + b << endl;
    return 0;
}`,
		"c": `#include <stdio.h>
int main() {
    int a, b;
    scanf("%d %d", &a, &b);
    printf("%d\\n", a + b);
    return 0;
}`,
		"java": `import java.util.Scanner;
public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        int a = sc.nextInt();
        int b = sc.nextInt();
        System.out.println(a + b);
    }
}`,
		"python": `a, b = map(int, input().split())
print(a + b)`,
	}

	submissions := []model.Submission{}
	subID := 0

	// Generate 55 submissions
	submissionDefs := []struct {
		userIdx    int
		problemIdx int
		verdict    string
		language   string
		daysAgo    int
	}{
		// AC submissions
		{1, 0, "AC", "cpp", 60}, {2, 0, "AC", "c", 55}, {3, 0, "AC", "java", 50},
		{4, 0, "AC", "python", 45}, {5, 0, "AC", "cpp", 40}, {1, 1, "AC", "cpp", 48},
		{2, 1, "AC", "python", 42}, {3, 1, "AC", "cpp", 38}, {4, 1, "AC", "java", 35},
		{5, 1, "AC", "cpp", 30}, {6, 0, "AC", "cpp", 28}, {6, 1, "AC", "cpp", 26},
		{7, 0, "AC", "python", 25}, {8, 0, "AC", "cpp", 22}, {9, 0, "AC", "c", 20},
		{1, 2, "AC", "cpp", 45}, {3, 2, "AC", "cpp", 40}, {5, 2, "AC", "python", 35},
		{1, 3, "AC", "cpp", 42}, {2, 3, "AC", "java", 38}, {4, 3, "AC", "cpp", 33},
		{6, 3, "AC", "cpp", 28}, {1, 4, "AC", "cpp", 40}, {3, 4, "AC", "cpp", 36},
		{5, 4, "AC", "python", 32}, {1, 5, "AC", "cpp", 38}, {2, 5, "AC", "cpp", 34},
		{4, 5, "AC", "java", 30}, {6, 5, "AC", "cpp", 26}, {7, 5, "AC", "python", 24},
		// WA submissions
		{2, 2, "WA", "cpp", 44}, {4, 2, "WA", "java", 39}, {7, 1, "WA", "python", 23},
		{8, 1, "WA", "cpp", 21}, {9, 1, "WA", "c", 19}, {3, 6, "WA", "cpp", 34},
		{5, 6, "WA", "python", 30}, {8, 5, "WA", "cpp", 23},
		// TLE submissions
		{1, 6, "TLE", "cpp", 36}, {2, 6, "TLE", "java", 33}, {4, 6, "TLE", "python", 31},
		{3, 7, "TLE", "cpp", 32}, {7, 6, "TLE", "python", 22},
		// RE submissions
		{6, 7, "RE", "cpp", 27}, {8, 7, "RE", "c", 20}, {9, 7, "RE", "python", 18},
		// CE submissions
		{9, 2, "CE", "cpp", 19}, {8, 2, "CE", "java", 21}, {7, 2, "CE", "python", 23},
		// MLE submissions
		{1, 8, "MLE", "cpp", 30}, {3, 8, "MLE", "cpp", 28}, {5, 8, "MLE", "python", 26},
		// More AC for harder problems
		{4, 10, "AC", "cpp", 25}, {6, 10, "AC", "cpp", 22}, {1, 15, "AC", "cpp", 20},
		{3, 15, "AC", "cpp", 18}, {5, 16, "AC", "cpp", 16}, {2, 19, "AC", "cpp", 15},
		{4, 19, "AC", "python", 14}, {6, 20, "AC", "cpp", 12},
	}

	for _, sd := range submissionDefs {
		if sd.userIdx >= len(users) || sd.problemIdx >= len(problems) {
			continue
		}
		code := ""
		if sd.verdict == "AC" {
			code = acCodeSnippets[sd.language]
		} else if sd.verdict == "CE" {
			code = "// compilation error\nint main( {"
		} else {
			code = acCodeSnippets[sd.language]
		}

		sub := model.Submission{
			UserID:       users[sd.userIdx].ID,
			ProblemID:    problems[sd.problemIdx].ID,
			Language:     sd.language,
			CodePath:     fmt.Sprintf("/submissions/%d/%d/code.%s", users[sd.userIdx].ID, subID, sd.language),
			CodeLength:   len(code),
			JudgeResult:  sd.verdict,
			JudgeScore:   0,
			TimeUsedMs:   0,
			MemoryUsedKb: 0,
			SubmitTime:   now.AddDate(0, 0, -sd.daysAgo),
		}
		if sd.verdict == "AC" {
			sub.JudgeScore = 100
			sub.TimeUsedMs = 10 + (subID % 50)
			sub.MemoryUsedKb = 1024 + (subID % 512)
		} else if sd.verdict == "TLE" {
			sub.TimeUsedMs = 2100
			sub.MemoryUsedKb = 2048
		} else if sd.verdict == "MLE" {
			sub.TimeUsedMs = 500
			sub.MemoryUsedKb = 300000
		} else if sd.verdict == "RE" {
			sub.ErrorMessage = "Runtime Error: segmentation fault"
		} else if sd.verdict == "WA" {
			sub.JudgeScore = 0
		} else if sd.verdict == "CE" {
			sub.ErrorMessage = "Compilation Error: expected ';'"
		}

		submissions = append(submissions, sub)
		subID++
	}

	for i := range submissions {
		db.FirstOrCreate(&submissions[i], model.Submission{
			UserID:    submissions[i].UserID,
			ProblemID: submissions[i].ProblemID,
			CodePath:  submissions[i].CodePath,
		})
	}

	// --- Contests ---
	contests := []model.Contest{
		{
			Title:        "YogduOJ 新手训练赛 #1",
			ContestType:  "individual",
			Category:     "programming",
			RuleType:     "acm",
			StartTime:    now.AddDate(0, 0, -30),
			EndTime:      now.AddDate(0, 0, -29),
			Description:  "欢迎参加 YogduOJ 新手训练赛！本次比赛包含 5 道入门级题目，适合刚接触算法竞赛的同学。\n\n## 比赛规则\n- ACM 赛制\n- 比赛时长 24 小时\n- 支持语言：C, C++, Java, Python",
			RuleDescription: "ACM赛制：通过题目数量优先，相同通过数则按罚时排名。每次提交错误会增加 20 分钟罚时。",
			MaxTeamSize:     1,
			AllowViewOthers: true,
			ShowRealtimeRank: true,
			Status:          "ended",
			ParticipantCount: 8,
			CreatedBy:       1,
			CreatedAt:       now.AddDate(0, -2, 0),
		},
		{
			Title:        "YogduOJ 周赛 #42",
			ContestType:  "individual",
			Category:     "programming",
			RuleType:     "acm",
			StartTime:    now.Add(-2 * time.Hour),
			EndTime:      now.Add(3 * time.Hour),
			Description:  "YogduOJ 第 42 场周赛，包含 5 道不同难度的题目。\n\n## 题目难度分布\n- A: 简单\n- B: 简单\n- C: 中等\n- D: 中等\n- E: 困难\n\n## 比赛规则\n- ACM 赛制\n- 比赛时长 5 小时\n- 支持语言：C, C++, Java, Python",
			RuleDescription: "ACM赛制：通过题目数量优先，相同通过数则按罚时排名。",
			MaxTeamSize:     1,
			AllowViewOthers: true,
			ShowRealtimeRank: true,
			Status:          "running",
			ParticipantCount: 6,
			CreatedBy:       1,
			CreatedAt:       now.AddDate(0, -1, 0),
		},
		{
			Title:        "YogduOJ 月赛 - 动态规划专题",
			ContestType:  "individual",
			Category:     "programming",
			RuleType:     "oi",
			StartTime:    now.AddDate(7, 0, 0),
			EndTime:      now.AddDate(7, 0, 0).Add(5 * time.Hour),
			Description:  "本次月赛聚焦动态规划专题，包含 4 道不同类型的 DP 题目。\n\n## 题目类型\n- 背包 DP\n- 区间 DP\n- 树形 DP\n- 状压 DP\n\n## 比赛规则\n- OI 赛制\n- 比赛时长 5 小时\n- 按总分排名，部分分可用",
			RuleDescription: "OI赛制：按总分排名，每道题根据通过的测试点给部分分。",
			MaxTeamSize:     1,
			AllowViewOthers: false,
			ShowRealtimeRank: false,
			Status:          "pending",
			ParticipantCount: 4,
			CreatedBy:       1,
			CreatedAt:       now,
		},
	}

	for i := range contests {
		db.FirstOrCreate(&contests[i], model.Contest{Title: contests[i].Title})
	}

	// --- Contest Problems ---
	contestProblems := []model.ContestProblem{
		// Contest 1 (ended) - problems 0,1,3,4,5
		{ContestID: contests[0].ID, ProblemID: problems[0].ID, DisplayOrder: 1, Score: 100, ProblemLabel: "A"},
		{ContestID: contests[0].ID, ProblemID: problems[1].ID, DisplayOrder: 2, Score: 100, ProblemLabel: "B"},
		{ContestID: contests[0].ID, ProblemID: problems[3].ID, DisplayOrder: 3, Score: 100, ProblemLabel: "C"},
		{ContestID: contests[0].ID, ProblemID: problems[4].ID, DisplayOrder: 4, Score: 100, ProblemLabel: "D"},
		{ContestID: contests[0].ID, ProblemID: problems[5].ID, DisplayOrder: 5, Score: 100, ProblemLabel: "E"},
		// Contest 2 (running) - problems 0,2,6,9,14
		{ContestID: contests[1].ID, ProblemID: problems[0].ID, DisplayOrder: 1, Score: 100, ProblemLabel: "A"},
		{ContestID: contests[1].ID, ProblemID: problems[2].ID, DisplayOrder: 2, Score: 100, ProblemLabel: "B"},
		{ContestID: contests[1].ID, ProblemID: problems[6].ID, DisplayOrder: 3, Score: 100, ProblemLabel: "C"},
		{ContestID: contests[1].ID, ProblemID: problems[9].ID, DisplayOrder: 4, Score: 100, ProblemLabel: "D"},
		{ContestID: contests[1].ID, ProblemID: problems[14].ID, DisplayOrder: 5, Score: 100, ProblemLabel: "E"},
		// Contest 3 (pending) - problems 5,6,9,12
		{ContestID: contests[2].ID, ProblemID: problems[5].ID, DisplayOrder: 1, Score: 25, ProblemLabel: "A"},
		{ContestID: contests[2].ID, ProblemID: problems[6].ID, DisplayOrder: 2, Score: 25, ProblemLabel: "B"},
		{ContestID: contests[2].ID, ProblemID: problems[9].ID, DisplayOrder: 3, Score: 25, ProblemLabel: "C"},
		{ContestID: contests[2].ID, ProblemID: problems[12].ID, DisplayOrder: 4, Score: 25, ProblemLabel: "D"},
	}

	for i := range contestProblems {
		db.FirstOrCreate(&contestProblems[i], model.ContestProblem{
			ContestID: contestProblems[i].ContestID,
			ProblemID: contestProblems[i].ProblemID,
		})
	}

	// --- Contest Signups ---
	contestSignups := []model.ContestSignup{
		// Contest 1 signups
		{ContestID: contests[0].ID, UserID: users[1].ID, SignupTime: now.AddDate(0, 0, -32)},
		{ContestID: contests[0].ID, UserID: users[2].ID, SignupTime: now.AddDate(0, 0, -31)},
		{ContestID: contests[0].ID, UserID: users[3].ID, SignupTime: now.AddDate(0, 0, -31)},
		{ContestID: contests[0].ID, UserID: users[4].ID, SignupTime: now.AddDate(0, 0, -30)},
		{ContestID: contests[0].ID, UserID: users[5].ID, SignupTime: now.AddDate(0, 0, -30)},
		{ContestID: contests[0].ID, UserID: users[6].ID, SignupTime: now.AddDate(0, 0, -30)},
		{ContestID: contests[0].ID, UserID: users[7].ID, SignupTime: now.AddDate(0, 0, -29)},
		{ContestID: contests[0].ID, UserID: users[8].ID, SignupTime: now.AddDate(0, 0, -29)},
		// Contest 2 signups
		{ContestID: contests[1].ID, UserID: users[1].ID, SignupTime: now.AddDate(0, -1, -2)},
		{ContestID: contests[1].ID, UserID: users[3].ID, SignupTime: now.AddDate(0, -1, -1)},
		{ContestID: contests[1].ID, UserID: users[4].ID, SignupTime: now.AddDate(0, -1, -1)},
		{ContestID: contests[1].ID, UserID: users[5].ID, SignupTime: now.AddDate(0, -1, 0)},
		{ContestID: contests[1].ID, UserID: users[7].ID, SignupTime: now.AddDate(0, -1, 0)},
		{ContestID: contests[1].ID, UserID: users[9].ID, SignupTime: now.AddDate(0, -1, 0)},
		// Contest 3 signups
		{ContestID: contests[2].ID, UserID: users[1].ID, SignupTime: now.AddDate(0, 0, 1)},
		{ContestID: contests[2].ID, UserID: users[3].ID, SignupTime: now.AddDate(0, 0, 1)},
		{ContestID: contests[2].ID, UserID: users[5].ID, SignupTime: now.AddDate(0, 0, 2)},
		{ContestID: contests[2].ID, UserID: users[7].ID, SignupTime: now.AddDate(0, 0, 2)},
	}

	for i := range contestSignups {
		db.FirstOrCreate(&contestSignups[i], model.ContestSignup{
			ContestID: contestSignups[i].ContestID,
			UserID:    contestSignups[i].UserID,
		})
	}

	// --- Announcements ---
	announcements := []model.Announcement{
		{
			Title:    "欢迎来到 YogduOJ！",
			Content:  "YogduOJ 正式上线啦！这是一个全新的在线评测平台，支持多种编程语言，提供丰富的算法题目和比赛功能。\n\n## 平台特色\n- 支持 C/C++、Java、Python 等多种语言\n- ACM/OI/IOI 等多种比赛模式\n- 实时排名和评测反馈\n- 团队协作功能\n\n如有问题请联系管理员 admin@yogdunana.com。",
			IsPinned:  true,
			CreatedBy: 1,
			CreatedAt: now.AddDate(-10, 0, 0),
		},
		{
			Title:    "YogduOJ 周赛 #42 正在进行中",
			Content:  "YogduOJ 第 42 场周赛正在进行中！比赛将于今天 23:00 结束，请还没有参加的同学抓紧时间。\n\n比赛链接：[周赛 #42](/contest/2)\n\n题目难度：A(简单) B(简单) C(中等) D(中等) E(困难)\n\n祝大家比赛顺利！",
			IsPinned:  true,
			CreatedBy: 1,
			CreatedAt: now.Add(-3 * time.Hour),
		},
		{
			Title:    "月赛预告：动态规划专题",
			Content:  "下周五将举办 YogduOJ 月赛 - 动态规划专题！\n\n## 比赛信息\n- 时间：下周六 14:00 - 19:00\n- 赛制：OI 赛制（部分分）\n- 题目数：4 道\n- 题目类型：背包 DP、区间 DP、树形 DP、状压 DP\n\n建议赛前复习相关知识点，做好准备！",
			IsPinned:  true,
			CreatedBy: 1,
			CreatedAt: now.AddDate(0, 0, -1),
		},
		{
			Title:    "系统维护通知",
			Content:  "为了提升平台性能和稳定性，我们计划于本周日凌晨 2:00 - 4:00 进行系统维护升级。\n\n维护期间平台将暂时无法访问，请提前做好准备。\n\n## 升级内容\n- 评测机性能优化\n- 新增 Python3.12 支持\n- 修复已知问题\n\n感谢大家的理解和支持！",
			IsPinned:  false,
			CreatedBy: 1,
			CreatedAt: now.AddDate(0, 0, -5),
		},
		{
			Title:    "新手训练赛 #1 赛后总结",
			Content:  "新手训练赛 #1 已圆满结束！感谢所有参赛选手的积极参与。\n\n## 比赛数据\n- 参赛人数：8 人\n- 完成全部题目的选手：3 人\n- 最受欢迎的题目：A+B Problem（通过率 95%）\n\n## 获奖名单\n- 第一名：sunqi（5 题全对，用时 1:23:45）\n- 第二名：wangwu（4 题对，罚时 2:10:30）\n- 第三名：zhangsan（4 题对，罚时 2:45:15）\n\n恭喜所有获奖选手！",
			IsPinned:  false,
			CreatedBy: 1,
			CreatedAt: now.AddDate(0, 0, -28),
		},
		{
			Title:    "新增题目集：图论入门",
			Content:  "我们新增了图论入门题目集，包含以下题目：\n\n1. 最小生成树\n2. 单源最短路\n3. 拓扑排序\n4. 图的着色\n\n建议按顺序练习，从易到难。祝大家学习愉快！",
			IsPinned:  false,
			CreatedBy: 1,
			CreatedAt: now.AddDate(0, 0, -7),
		},
	}

	for i := range announcements {
		db.FirstOrCreate(&announcements[i], model.Announcement{Title: announcements[i].Title})
	}

	// --- Teams ---
	teams := []model.Team{
		{
			Name:       "算法之星",
			Slogan:     "用算法改变世界",
			LeaderID:   users[5].ID, // sunqi
			MaxMembers: 3,
			Status:     "active",
			CreatedAt:  now.AddDate(-5, 0, 0),
		},
		{
			Name:       "代码猎人",
			Slogan:     "Bug 猎手，代码无错",
			LeaderID:   users[3].ID, // wangwu
			MaxMembers: 3,
			Status:     "active",
			CreatedAt:  now.AddDate(-4, 0, 0),
		},
		{
			Name:       "ACMer联盟",
			Slogan:     "AC 永不为奴",
			LeaderID:   users[1].ID, // zhangsan
			MaxMembers: 3,
			Status:     "active",
			CreatedAt:  now.AddDate(-3, 0, 0),
		},
	}

	for i := range teams {
		db.FirstOrCreate(&teams[i], model.Team{Name: teams[i].Name})
	}

	// --- Team Members ---
	teamMembers := []model.TeamMember{
		// Team 1: 算法之星 (leader: sunqi)
		{TeamID: teams[0].ID, UserID: users[5].ID, Role: "leader", JoinedAt: now.AddDate(-5, 0, 0)},
		{TeamID: teams[0].ID, UserID: users[3].ID, Role: "member", JoinedAt: now.AddDate(-5, 0, 1)},
		{TeamID: teams[0].ID, UserID: users[7].ID, Role: "member", JoinedAt: now.AddDate(-4, 0, 0)},
		// Team 2: 代码猎人 (leader: wangwu)
		{TeamID: teams[1].ID, UserID: users[3].ID, Role: "leader", JoinedAt: now.AddDate(-4, 0, 0)},
		{TeamID: teams[1].ID, UserID: users[1].ID, Role: "member", JoinedAt: now.AddDate(-4, 0, 2)},
		{TeamID: teams[1].ID, UserID: users[9].ID, Role: "member", JoinedAt: now.AddDate(-3, 0, 0)},
		// Team 3: ACMer联盟 (leader: zhangsan)
		{TeamID: teams[2].ID, UserID: users[1].ID, Role: "leader", JoinedAt: now.AddDate(-3, 0, 0)},
		{TeamID: teams[2].ID, UserID: users[2].ID, Role: "member", JoinedAt: now.AddDate(-3, 0, 1)},
		{TeamID: teams[2].ID, UserID: users[4].ID, Role: "member", JoinedAt: now.AddDate(-2, 0, 0)},
	}

	for i := range teamMembers {
		db.FirstOrCreate(&teamMembers[i], model.TeamMember{
			TeamID: teamMembers[i].TeamID,
			UserID: teamMembers[i].UserID,
		})
	}

	// --- Team Invitations (pending) ---
	teamInvitations := []model.TeamInvitation{
		{
			TeamID:    teams[0].ID,
			InviterID: users[5].ID,
			InviteeID: users[4].ID,
			Status:    "pending",
			CreatedAt: now.AddDate(0, 0, -2),
		},
		{
			TeamID:    teams[1].ID,
			InviterID: users[3].ID,
			InviteeID: users[6].ID,
			Status:    "pending",
			CreatedAt: now.AddDate(0, 0, -3),
		},
		{
			TeamID:    teams[2].ID,
			InviterID: users[1].ID,
			InviteeID: users[8].ID,
			Status:    "pending",
			CreatedAt: now.AddDate(0, 0, -1),
		},
	}

	for i := range teamInvitations {
		db.FirstOrCreate(&teamInvitations[i], model.TeamInvitation{
			TeamID:    teamInvitations[i].TeamID,
			InviterID: teamInvitations[i].InviterID,
			InviteeID: teamInvitations[i].InviteeID,
		})
	}

	// --- System Configs ---
	configs := []model.SystemConfig{
		{ConfigKey: "site_name", ConfigValue: "YogduOJ", Description: "站点名称"},
		{ConfigKey: "site_description", ConfigValue: "YogduOJ - 在线评测平台", Description: "站点描述"},
		{ConfigKey: "default_time_limit", ConfigValue: "2000", Description: "默认时间限制（毫秒）"},
		{ConfigKey: "default_memory_limit", ConfigValue: "262144", Description: "默认内存限制（KB）"},
		{ConfigKey: "supported_languages", ConfigValue: `["cpp","c","java","python"]`, Description: "支持的编程语言"},
		{ConfigKey: "max_code_length", ConfigValue: "65536", Description: "代码最大长度（字节）"},
		{ConfigKey: "register_enabled", ConfigValue: "true", Description: "是否开放注册"},
	}

	for i := range configs {
		db.FirstOrCreate(&configs[i], model.SystemConfig{ConfigKey: configs[i].ConfigKey})
	}

	return nil
}
