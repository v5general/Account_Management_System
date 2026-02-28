<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  const token = localStorage.getItem('token')
  console.log('Token exists:', !!token)
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  width: 100%;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

/* 修复 Element Plus el-select 选中值不显示的问题 */
.el-select__wrapper {
  display: inline-flex !important;
  align-items: center !important;
}

.el-select__selection {
  display: flex !important;
  flex-wrap: wrap !important;
  position: relative !important;
}

.el-select__selected-item {
  display: inline-flex !important;
  align-items: center !important;
  visibility: visible !important;
  opacity: 1 !important;
  color: var(--el-text-color-regular) !important;
  font-size: 14px !important;
}

.el-select__selected-item span {
  display: inline !important;
  visibility: visible !important;
  color: var(--el-text-color-regular) !important;
}

.el-select__input {
  color: var(--el-text-color-regular) !important;
  opacity: 1 !important;
}

.el-select__placeholder {
  color: var(--el-text-color-placeholder) !important;
}

/* 强制显示选中值的文本 */
.el-select .el-select__wrapper .el-select__selection .el-select__selected-item {
  color: #606266 !important;
}

.el-select .el-select__wrapper .el-select__selection .el-select__selected-item span {
  color: #606266 !important;
}
</style>
