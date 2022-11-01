<script setup lang="ts">
// This template is using Vue 3 <script setup> SFCs
// Check out https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup
import {reactive, ref, provide, InjectionKey, ComponentPublicInstance} from 'vue'
import { ElMessage } from 'element-plus'
// import {
//   loginValidate,
// } from '@/api/common'

import {
  SaveUser, CheckUser
} from '@/../wailsjs/go/app/api'
import { app } from '@/../wailsjs/go/models'

interface Emits {
  (e: 'login-success'): void
}
const emit = defineEmits<Emits>()

const loginFormFef = ref()

const form = reactive({
  username: "",
  password: ""
})

const rules = {
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: 'blur',
    },
    {
      min: 4,
      max: 30,
      message: '请输入正确的用户名',
      trigger: 'blur',
    }
  ],
  password: [
    {
      required: true,
      message: '请输入授权码',
      trigger: 'blur',
    },
    {
      min: 4,
      max: 50,
      message: '请输入正确的授权码',
      trigger: 'blur',
    }
  ]
}

const submitForm = ()=> {
  loginFormFef.value.validate(async(v:boolean) =>{
    if (v) {
      const userSecrect = app.userSecret.createFrom({
        username: form.username,
        key: form.password,
      })
      SaveUser(userSecrect).then(res => {
        console.log("Validate")
        console.log(res)
        if (res) {
          emit('login-success')
          ElMessage({
            type: 'success',
            message: '校验成功，欢迎使用',
          })
        } else {
          ElMessage({
            type: 'error',
            message: '请填写正确的用户名或授权码',
            showClose: true
          })
          form.username = ""
          form.password = ""
          
        }
      }).catch(err => {
        console.log(err)
        ElMessage({
          type: 'error',
          message: '请填写正确的用户名或授权码',
          showClose: true
        })
      })

    } else {
      ElMessage({
        type: 'error',
        message: '请填写正确的用户名或授权码',
        showClose: true
      })

    }
  })
}
</script>

<template>
  <el-form ref="loginFormFef" :model="form" :rules="rules" label-width="70px">
  <el-form-item label="用户名" prop="username">
    <el-input v-model="form.username" />
  </el-form-item>
  <el-form-item label="授权码" prop="password">
    <el-input v-model="form.password" />
  </el-form-item>
  <el-form-item>
    <el-button
      type="primary"
      style="width: 100%; margin-left: 0%"
      @click="submitForm"
      >提交
    </el-button>
  </el-form-item>

  </el-form>

</template>


<style>

</style>
