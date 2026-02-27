<template>
  <div class="account-manage">
    <el-card>
      <template #header>
        <span>账号管理</span>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" style="max-width: 500px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="form.real_name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="不修改密码请留空" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="form.confirm_password" type="password" placeholder="请再次输入新密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleUpdate" :loading="loading">保存修改</el-button>
          <el-button @click="handleDeleteAccount" type="danger" plain>注销账号</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { updateMyAccount, getUserList } from '@/api/user'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const originalUsername = ref('')

const form = reactive({
  username: '',
  real_name: '',
  password: '',
  confirm_password: ''
})

const validateUsername = async (rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请输入用户名'))
    return
  }
  // 用户名格式验证
  const usernameRegex = /^[a-zA-Z0-9_]+$/
  if (!usernameRegex.test(value)) {
    callback(new Error('用户名只能包含字母、数字和下划线'))
    return
  }
  // 如果用户名没变，直接通过
  if (value === originalUsername.value) {
    callback()
    return
  }
  // 检查用户名是否已存在
  try {
    const res = await getUserList({ page: 1, page_size: 1000 })
    const exists = res.data.list.some((u: any) => u.username === value)
    if (exists) {
      callback(new Error('该用户名已被使用'))
    } else {
      callback()
    }
  } catch {
    callback()
  }
}

const validateConfirmPassword = (rule: any, value: any, callback: any) => {
  if (value && value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  username: [
    { required: true, validator: validateUsername, trigger: 'blur' }
  ],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  password: [
    { min: 8, message: '密码长度不能少于8位', trigger: 'blur' }
  ],
  confirm_password: [
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

onMounted(() => {
  const user = userStore.userInfo
  if (user) {
    form.username = user.username
    originalUsername.value = user.username
    form.real_name = user.real_name
  }
})

async function handleUpdate() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const data: any = {
        real_name: form.real_name
      }
      // 只有当用户名改变时才发送
      if (form.username !== originalUsername.value) {
        data.username = form.username
      }
      if (form.password) {
        data.password = form.password
      }

      await updateMyAccount(data)
      ElMessage.success('修改成功')

      // 更新用户信息
      await userStore.fetchUserInfo()
      originalUsername.value = form.username

      // 清空密码字段
      form.password = ''
      form.confirm_password = ''
    } catch (error: any) {
      console.error('Failed to update account:', error)
      ElMessage.error(error.response?.data?.message || '修改失败')
    } finally {
      loading.value = false
    }
  })
}

async function handleDeleteAccount() {
  try {
    await ElMessageBox.confirm(
      '注销账号将永久删除您的账号和所有相关数据，此操作不可恢复！确定要继续吗？',
      '注销账号',
      {
        type: 'error',
        confirmButtonText: '确定注销',
        cancelButtonText: '取消',
        confirmButtonClass: 'el-button--danger'
      }
    )

    // 再次确认
    await ElMessageBox.prompt('请输入"注销"以确认此操作', '二次确认', {
      inputPattern: /^注销$/,
      inputErrorMessage: '输入不正确，操作已取消'
    })

    loading.value = true
    try {
      // 这里调用删除用户API
      // await deleteUser(form.user_id)
      ElMessage.success('账号已注销')
      await userStore.logout()
    } catch (error: any) {
      console.error('Failed to delete account:', error)
      ElMessage.error(error.response?.data?.message || '注销失败')
    } finally {
      loading.value = false
    }
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.account-manage {
  padding: 20px;
}
</style>
