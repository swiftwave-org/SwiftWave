<script setup>
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api'
import 'monaco-editor/esm/vs/basic-languages/dockerfile/dockerfile.contribution.js'

import { onMounted, ref } from 'vue'

const props = defineProps({
  code: {
    type: String,
    default: ''
  },
  codeChanged: {
    type: Function,
    default: () => {}
  }
})

const editor = ref()

onMounted(() => {
  const editorInstance = monaco.editor.create(editor.value, {
    value: props.code,
    language: 'dockerfile'
  })

  editor.value.addEventListener('keyup', () => {
    props.codeChanged(editorInstance.getValue())
  })
})
</script>

<template>
  <div id="editor" ref="editor"></div>
</template>

<style scoped>
#editor {
  width: 100%;
  height: 75vh;
}
</style>
