<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import type * as Monaco from 'monaco-editor'

const props = withDefaults(defineProps<{
  modelValue: string
  language?: string
  readonly?: boolean
  height?: string
}>(), {
  modelValue: '',
  language: 'cpp',
  readonly: false,
  height: '400px',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const themeStore = useThemeStore()
const editorRef = ref<HTMLElement | null>(null)
const monacoEditor = ref<Monaco.editor.IStandaloneCodeEditor | null>(null)
const monacoInstance = ref<typeof Monaco | null>(null)
const isReady = ref(false)

const editorTheme = computed(() => themeStore.isDark ? 'vs-dark' : 'vs')

const languageMap: Record<string, string> = {
  c: 'c',
  cpp: 'cpp',
  java: 'java',
  python3: 'python',
  python: 'python',
  javascript: 'javascript',
  go: 'go',
  rust: 'rust',
}

const mappedLanguage = computed(() => languageMap[props.language] || props.language)

async function initEditor() {
  try {
    const monaco = await import('monaco-editor')
    monacoInstance.value = monaco

    monacoEditor.value = monaco.editor.create(editorRef.value!, {
      value: props.modelValue,
      language: mappedLanguage.value,
      theme: editorTheme.value,
      readOnly: props.readonly,
      minimap: { enabled: false },
      lineNumbers: 'on',
      wordWrap: 'on',
      autoIndent: 'full',
      formatOnPaste: true,
      formatOnType: true,
      scrollBeyondLastLine: false,
      fontSize: 14,
      tabSize: 2,
      renderWhitespace: 'selection',
      automaticLayout: true,
      scrollbar: {
        verticalScrollbarSize: 8,
        horizontalScrollbarSize: 8,
      },
    })

    monacoEditor.value.onDidChangeModelContent(() => {
      const value = monacoEditor.value!.getValue()
      emit('update:modelValue', value)
    })

    isReady.value = true
  } catch (e) {
    console.error('Failed to load Monaco Editor:', e)
  }
}

watch(() => props.language, (newLang) => {
  if (monacoEditor.value && monacoInstance.value) {
    const lang = languageMap[newLang] || newLang
    monacoInstance.value.editor.setModelLanguage(monacoEditor.value.getModel()!, lang)
  }
})

watch(() => props.modelValue, (newValue) => {
  if (monacoEditor.value && monacoEditor.value.getValue() !== newValue) {
    monacoEditor.value.setValue(newValue)
  }
})

watch(editorTheme, (newTheme) => {
  if (monacoInstance.value) {
    monacoInstance.value.editor.setTheme(newTheme)
  }
})

watch(() => props.readonly, (newReadonly) => {
  if (monacoEditor.value) {
    monacoEditor.value.updateOptions({ readOnly: newReadonly })
  }
})

onMounted(() => {
  initEditor()
})
</script>

<template>
  <div class="code-editor-wrapper">
    <div v-if="!isReady" class="code-editor-loading">
      <NSpin size="small" />
      <span>{{ $t('common.loading') }}</span>
    </div>
    <div ref="editorRef" class="code-editor-container" :style="{ height }"></div>
  </div>
</template>

<style scoped lang="scss">
.code-editor-wrapper {
  position: relative;
  width: 100%;
  border-radius: var(--border-radius);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.code-editor-loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background-color: var(--color-bg-secondary);
  color: var(--color-text-secondary);
  font-size: 14px;
  z-index: 1;
}

.code-editor-container {
  width: 100%;
  min-height: 200px;
}
</style>
