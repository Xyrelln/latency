<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import Automation from '../components/Automation.vue';
import AboutPage from '../components/AboutPage.vue';

import AutomationPage from './AutomationPage.vue';
import AndroidLatencyPage from './AndroidLatencyPage.vue';
import PCLatencyPage from './PCLatencyPage.vue';
import Login from '@/components/Login.vue'

import { EventsOn, EventsOff } from '@/../wailsjs/runtime/runtime'
import { CheckUpdate, DoUpdate, CheckUser, SaveUser } from '@/../wailsjs/go/app/api'
import { app } from '@/../wailsjs/go/models'
import { isWailsRun } from '@/utils/utils'


const topTabName = ref('latency-pc')
const upgradeDialogVisible = ref(false)
const upgradeInfo = reactive({
  latestVersion: '',
  needUpdate: false,
})
const isLoginFormShow = ref(false)

/**
 * 开启事件监听
 */
const setEventOn = () => {
  EventsOn("latency:update_success", (data: any) => {
    ElMessage({
      type: 'success',
      message: '新版本更新成功，下次打开应用生效',
      showClose: true
    })
  })

  EventsOn("latency:update_error", (err: any) => {
    ElMessage({
      type: 'error',
      message: '新版本更新失败: ' + err,
      showClose: true
    })
  })
}

/**
 * 关闭事件监听
 */
const setEventOff = () => {
  EventsOff('latency:update_success')
  EventsOff('latency:update_error')
}

/**
 * 升级检测
 */
 const handleCheckUpgrade = async() => {
  CheckUpdate().then( (res:app.UpdateInfo) => {
    // console.log(res)
    if (res.needUpdate === true) {
      upgradeInfo.needUpdate = res.needUpdate
      upgradeInfo.latestVersion = res.latestVersion

      upgradeDialogVisible.value = true
    }
  }).catch(err => { 
    console.log(err)
  })
}

const handleCheckUser = () => {
  CheckUser().then(res => {
    console.log(res)
  }).catch(err => {
    console.log(err)
    isLoginFormShow.value = true
  })
}

/**
 * 执行升级
 */
 const handleUpgrade = async() => {
  DoUpdate(upgradeInfo.latestVersion).then(res => {

  }).catch(err => {
    console.log(err)
  })
  
  upgradeDialogVisible.value = false
}

const handleLoginSuccess = ()=> {
  isLoginFormShow.value = false
}

const handleCloseLogin = () => {
  isLoginFormShow.value = true
}

onMounted(()=>{
  if (isWailsRun()) {
    setEventOn()
    handleCheckUpgrade()
    handleCheckUser()
  } 
})

onUnmounted(()=>{
  if (isWailsRun()) {
    setEventOff()
  } 
})


</script>

<template>
  <el-container class="main-container">
    <el-tabs type="border-card" v-model="topTabName" class="top-tabs" >
      <el-tab-pane label="延迟测试-PC" name="latency-pc" class="latency-panel">
        <PCLatencyPage></PCLatencyPage>
      </el-tab-pane>
      <el-tab-pane label="延迟测试-Android" name="latency" class="latency-panel">
        <AndroidLatencyPage></AndroidLatencyPage>
      </el-tab-pane>
      <el-tab-pane label="场景配置-Android" name="automation" class="latency-panel">
        <AutomationPage></AutomationPage>
      </el-tab-pane>
      <el-tab-pane label="关于" name="about" class="latency-panel">
        <AboutPage></AboutPage>
      </el-tab-pane>
    </el-tabs>

    <el-dialog
      v-model="upgradeDialogVisible"
      title="升级提示"
      width="30%"
    >
      <span>检测到新版本: {{ upgradeInfo.latestVersion }} , 是否自动升级</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="upgradeDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUpgrade"
            >确认升级</el-button
          >
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="isLoginFormShow" title="授权校验" width="35%" center :before-close="handleCloseLogin">
      <Login 
        ref="loginRef"
        @login-success="handleLoginSuccess"
      />
    </el-dialog>
   
  </el-container>
    
</template>

<style>
.top-tabs {
  width: 100vw;
  /* width: calc(100vw - 14px); */
  /* height: inherit; */
}
.latency-panel {
  height: inherit;
}

/* .automation-panel {
  height: 80vh;
} */

.el-tabs--border-card>.el-tabs__content {
  padding: 7px !important;
  height: calc(100vh - 74px);
}

.main-container {
  width: 100vw;
  height: 100vh;
}
</style>
