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
</style>

<style>
/* 修复 Element Plus el-select 选中值不显示的问题 - 使用更高优先级 */
.el-select__wrapper .el-select__selection .el-select__selected-item.el-select__input-wrapper.is-hidden {
  display: inline-flex !important;
  visibility: visible !important;
  opacity: 1 !important;
  position: static !important;
}

.el-select__wrapper .el-select__selection .el-select__selected-item {
  display: inline-flex !important;
  align-items: center !important;
  visibility: visible !important;
  opacity: 1 !important;
  color: #606266 !important;
  font-size: 14px !important;
}

.el-select__wrapper .el-select__selection .el-select__selected-item span {
  display: inline !important;
  visibility: visible !important;
  color: #606266 !important;
}
</style>
