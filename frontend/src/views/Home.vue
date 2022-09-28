<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { UserFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import ImagePreview from '../components/ImagePreview.vue';
import FileRecord from '../components/FileRecord.vue';
import Automation from '../components/Automation.vue';
import AboutPage from '../components/AboutPage.vue';
import HelpPage from '../components/HelpPage.vue';

import { 
  ListDevices,
  Start,
  StartRecord,
  // StopRecord,
  // StopProcessing,
  StartTransform,
  StartAnalyse,
  SetPointerLocationOff,
  SetPointerLocationOn,
  GetFirstImageInfo,
  ClearCacheData,
  SetAutoSwipeOn,
  SetAutoSwipeOff,
  GetDisplay,
  GetImageFiles,
  InputSwipe,
  IsAppReady,
  StartWithVideo,
  GetPhysicalSize,
  ListRecords,
  // IsAppReady2,
} from '../../wailsjs/go/app/Api'
import {adb, core} from '../../wailsjs/go/models'
import {
  EventsOn,
  EventsOff,
  WindowReload,
} from '../../wailsjs/runtime/runtime'

const deviceSelected = ref("")
const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})

const topTabName = ref('latency')
const latencyTabName = ref('list')
const placeholder = "./src/assets/images/placeholder.png"

const fileRecordRef = ref()

const status = ref(false)
const form = reactive({
  device: '',
  sx: 0,
  sy: 0,
  dx: 0,
  dy: 0,
  speed: 500,
  interval: 2000,
  scene_id: '',
  location: true
})

const interval = ref()
const processStatus = ref(0)
const countDownSecond = ref(0)
const imagePreviewRef = ref()
const externalVideoPath = ref('')
const deviceInfo = reactive({
  width: 1080,
  height: 1920,
})

const imageInfo = reactive({
  path: placeholder,
  width: 0,
  height: 0,
  size: 0
})

const settingForm = reactive({
  touchScore: 4,
  diffScore: 20,
  timeout: 3,
  sceneStart: 1500,
  prepareTimeout: 3,
  develop: false,
  autoUpload: false
})

provide('threshold', settingForm.diffScore)


const rules =  {
  touchScore: [
    {
      required: true,
      message: '触控比对阈值',
      trigger: 'blur',
    },
    {
      validator: checkGreaterThanZero,
      trigger: 'blur',
    }
  ],
  diffScore: [
    {
      required: true,
      message: '选中区域比对阈值',
      trigger: 'blur',
    },
    {
      validator: checkGreaterThanZero,
      trigger: 'blur',
    }
  ],
  timeout: [
    {
      required: true,
      message: '录制时长',
      trigger: 'blur',
    },
    {
      validator: checkGreaterThanZero,
      trigger: 'blur',
    }
  ],
  sceneStart: [
    {
      required: true,
      message: '场景相对录制开始时间',
      trigger: 'blur',
    },
    {
      validator: checkGreaterThanZero,
      trigger: 'blur',
    }
  ],
}

const scenes = ref([
  {
    name: "第五人格-视角移动",
    value: "1",
    behavior: {
      type: 'swipe',
      location: {
        sx: deviceInfo.height/2,
        sy:  deviceInfo.width/2,
        dx: deviceInfo.height/2 + deviceInfo.height/2/2,
        dy: deviceInfo.width/2,
        speed: form.speed
      }
    },
    monitor_rect: {
      top: 26,
      left: 0,
      width: 466,
      height: 90
    },
  }
])


function checkGreaterThanZero (rule: any, value: any, callback: any)  {
  if (value <= 0) {
    callback(new Error('数值必须大于0'))
  } else {
    callback()
  }
}


/**
 * 获取设备列表
 * @param value 设备序列号
 */
function getDeviceList (value: any) {``
  ListDevices().then(result => {
    data.devices = result
  }).catch(err => {
    ElMessage({
      type: 'error',
      message: '设备获取失败, error: ' + err
    })
  })
}

/**
 * 获取选中设备状态
 */
function getSelectDeviceState(){
  for(let d of data.devices) {
    if (d.Serial == deviceSelected.value) {
      return d.State
    }
  }
}

function handleLoadExtVideo() {
  handleResetStatus()
  NProgress.start()
  StartWithVideo(externalVideoPath.value)
}

/**
 * 发送拖动事件
 */
function handleInputSwipe() {
  const swipeEvent = adb.SwipeEvent.createFrom(
    { 
      sx: deviceInfo.height/2,
      sy:  deviceInfo.width/2,
      dx: deviceInfo.height/2 + deviceInfo.height/2/2,
      dy: deviceInfo.width/2,
      speed: form.speed
    }
  )
  // const interval = 2
  console.log(swipeEvent)
  InputSwipe(deviceSelected.value, swipeEvent)
}

/**
 * 启动录制
 */
async function handleStartRecord() {
  // 设备状态检查
  const state = getSelectDeviceState()
  if (state == 0) {
    ElMessage({
      type: 'error',
      message: '设备已离线，请检查设备'
    })
    return
  }
  else if (state == 2) {
    ElMessage({
      type: 'error',
      message: '设备未授权，请检查设备'
    })
    return
  }

  handleResetStatus()
  NProgress.start()
  const result = await setPointerLocationOn()

  const status = await handleGetDisplay()
  if (!status) {
    await handleGetPhysicalSize()
  }

  Start(deviceSelected.value, settingForm.timeout)

  // 1s 后拖动
  // doSwipe()

  // add stop count down
  processStatus.value = 2
  runUntilCountDown(settingForm.timeout)

  setTimeout(() => {
    handleInputSwipe()
  }, 1500);
}

function handleStart() {
  StartRecord(deviceSelected.value)
}


function handleStopRecord() {
}

function handleToImage() {
  StartTransform()
}

function clearCurrentInterval() {
  if (interval.value != null) {
    clearInterval(interval.value)
    interval.value = null
  }
}

function runUntilCountDown(second: number, callback?: Function){
  countDownSecond.value = second
  function countDown() {
    // 启动需要时间，提前1s启动
    if (countDownSecond.value  == 1) {
      if (callback) { callback() }
    }
    if (countDownSecond.value  > 0) {
      countDownSecond.value  --
    }

  }

  clearCurrentInterval()
  interval.value = setInterval(countDown, 1000)
}

/**
 * 配置默认状态
 */
function handleResetStatus() {
  if (NProgress.isStarted()) {
    NProgress.done()
  }
  imagePreviewRef.value.setCalcButtonDisable(true)
  imagePreviewRef.value.setImagePlaceHolder()
  imagePreviewRef.value.setDefaultTime()
}


/**
 * 开启指针显示
 */
async function setPointerLocationOn():Promise<Boolean> {
  let result = false
  await SetPointerLocationOn(deviceSelected.value).then(res =>{ 
    ElMessage({
      type: 'success',
      message: '开启指针成功'
    })
     result = true
    }
  ).catch(err => {
    ElMessage({
      type: 'error',
      message: '开启指针失败'
    })
    result = false
  })
  return result
}

function setPointerLocationOff():Boolean {
  SetPointerLocationOff(deviceSelected.value).then(res =>{ 
      ElMessage({
        type: 'success',
        message: '关闭指针成功'
      })
      return true
  }).catch(err => {
    ElMessage({
      type: 'error',
      message: '关闭指针失败'
    })
  })
  return false
}

/**
 * 绑定监听
 */
async function addEventLister() {
  EventsOn("latency:record_start", ()=>{
    console.log("record_start")
    ElNotification({
      title: '进度提示: 1/3',
      type: 'info',
      message: "开始录制",
    })
  })
  EventsOn("latency:record_filish", ()=>{
    setPointerLocationOff()
    processStatus.value = 0
    ElNotification({
      title: '进度提示: 1/3',
      type: 'success',
      message: "录制成功",
    })
  })
  EventsOn("latency:transform_start", ()=>{
    ElNotification({
      title: '进度提示: 2/3',
      type: 'success',
      message: "数据预处理中，请稍后...",
      duration: 6000,
    })
    NProgress.done()
  })
  EventsOn("latency:transform_start_error", ()=>{
    ElNotification({
      title: '进度提示: 2/3',
      type: 'error',
      message: "数据预处理失败，请重试",
      duration: 0,
    })
    NProgress.done()
  })
  EventsOn("latency:record_start_error", ()=>{
    ElNotification({
      title: '进度提示: 2/3',
      type: 'error',
      message: "录制失败，请重试",
      duration: 0,
    })
    NProgress.done()
  })
  EventsOn("latency:transform_filish", ()=>{
    ElNotification({
      title: '进度提示: 2/3',
      type: 'success',
      message: "数据预处理完成，加载首帧画面",
    })
    getFirstImage()
    NProgress.done()
  })
}

function getFirstImage(){
  GetFirstImageInfo().then((res: core.ImageInfo) => {
    imageInfo.path = res.path
    imageInfo.width = res.width
    imageInfo.height = res.height
    imagePreviewRef.value.loadNewImage(res)
    imagePreviewRef.value.enableCalcButton()
  })
}

async function handleGetPhysicalSize() {
  let status = false
  await GetPhysicalSize(deviceSelected.value).then((res: adb.Display) => {
      deviceInfo.width = res.width
      deviceInfo.height = res.height
  }).catch(err => {
    // deviceInfo.width = 1080
    // deviceInfo.height = 1920
    console.log(err)
  })
  return status
}


async function handleGetDisplay() {
  let status = false
  await GetDisplay(deviceSelected.value).then((res: adb.Display) => {
      deviceInfo.width = res.width
      deviceInfo.height = res.height
      status = true
  }).catch(err => {
    // deviceInfo.width = 1080
    // deviceInfo.height = 1920
    console.log(err)
  })
  return status
}

async function removeEventLister() {
  EventsOff("latency:record_start")
  EventsOff("latency:record_filish")
  EventsOff("latency:transform_start")
  EventsOff("latency:transform_start_error")
  EventsOff("latency:record_start_error")
  EventsOff("latency:transform_filish")
}

/**
 * 环境检查
 */
async function initCheck() {
  IsAppReady().then(res => {
  }).catch(err => {
    ElNotification({
      title: '环境检查',
      type: 'error',
      message: err,
      duration: 0,
    })
  })
}

/**
 * 是否运行在 wails 环境
 */
function isWailsRun(){
  const tag = 'wails.io'
  return navigator.userAgent.indexOf(tag) != -1
}

onMounted(()=> {
  // 如果是在 wails 运行环境则运行环境检查及事件监听
  if (isWailsRun()) {
    initCheck()
    addEventLister()
    fileRecordRef.value.handleLoadCacheFiles()
  }
})

function handleClearCache() {
  ClearCacheData()
}

function handleReload() {
  WindowReload();
}

function handleStopProcessing() {}

onUnmounted(()=>{
  // removeEventLister()
  if (isWailsRun()) {
    removeEventLister()
  }
})

function handleGetImage() {
  imagePreviewRef.value.handleGetImage()
}



</script>

<template>
  <el-container>
    <el-tabs type="border-card" v-model="topTabName" >
      <el-tab-pane label="延迟测试" name="latency">
        <el-scrollbar style="height: calc(100vh - 100px);width: calc(100vw - 60px)">
          <el-container>
            <el-aside class="aside-content" width="240px">
              <el-row class="row-item">
                <el-col :span="20">
                  <el-select
                    v-model="deviceSelected"
                    @focus="getDeviceList"
                    filterable
                    placeholder="请选择设备"
                    style="width:100%">
                    <el-option
                      v-for="item in data.devices"
                      :key="item.Serial"
                      :label="item.Serial"
                      :value="item.Serial"
                    >
                    </el-option>
                  </el-select>

                </el-col>
                <el-col :span="4">
                  <el-tooltip
                    class="device-question"
                    effect="dark"
                    content="如列表为空，请检查设备是否正常连接"
                    placement="right"
                  >
                    <i class="el-icon button-icon" style="float: right;">
                      <svg t="1663058405930" class="icon button-icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3677" width="200" height="200"><path d="M512 784.352m-48 0a1.5 1.5 0 1 0 96 0 1.5 1.5 0 1 0-96 0Z" p-id="3678" fill="#8a8a8a"></path><path d="M512 960C264.96 960 64 759.04 64 512S264.96 64 512 64s448 200.96 448 448S759.04 960 512 960zM512 128.288C300.416 128.288 128.288 300.416 128.288 512c0 211.552 172.128 383.712 383.712 383.712 211.552 0 383.712-172.16 383.712-383.712C895.712 300.416 723.552 128.288 512 128.288z" p-id="3679" fill="#8a8a8a"></path><path d="M512 673.696c-17.664 0-32-14.336-32-32l0-54.112c0-52.352 40-92.352 75.328-127.648C581.216 434.016 608 407.264 608 385.92c0-53.344-43.072-96.736-96-96.736-53.824 0-96 41.536-96 94.56 0 17.664-14.336 32-32 32s-32-14.336-32-32c0-87.424 71.776-158.56 160-158.56s160 72.096 160 160.736c0 47.904-36.32 84.192-71.424 119.296C572.736 532.992 544 561.728 544 587.552l0 54.112C544 659.328 529.664 673.696 512 673.696z" p-id="3680" fill="#8a8a8a"></path></svg>
                    </i>
                  </el-tooltip>
                </el-col>
              </el-row>
              <el-row class="row-item">
                <el-col :span="20">
                  <el-select
                    v-model="deviceSelected"
                    filterable
                    placeholder="请选择场景"
                    style="width:100%">
                    <el-option
                      v-for="item in scenes"
                      :key="item.value"
                      :label="item.name"
                      :value="item.value"
                    >
                    </el-option>
                  </el-select>
                </el-col>
                <el-col :span="4">
                  <el-tooltip
                    class="device-question"
                    effect="dark"
                    content="选择自动运行场景"
                    placement="right"
                  >
                    <i class="el-icon button-icon" style="float: right;">
                      <svg t="1663058405930" class="icon button-icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3677" width="200" height="200"><path d="M512 784.352m-48 0a1.5 1.5 0 1 0 96 0 1.5 1.5 0 1 0-96 0Z" p-id="3678" fill="#8a8a8a"></path><path d="M512 960C264.96 960 64 759.04 64 512S264.96 64 512 64s448 200.96 448 448S759.04 960 512 960zM512 128.288C300.416 128.288 128.288 300.416 128.288 512c0 211.552 172.128 383.712 383.712 383.712 211.552 0 383.712-172.16 383.712-383.712C895.712 300.416 723.552 128.288 512 128.288z" p-id="3679" fill="#8a8a8a"></path><path d="M512 673.696c-17.664 0-32-14.336-32-32l0-54.112c0-52.352 40-92.352 75.328-127.648C581.216 434.016 608 407.264 608 385.92c0-53.344-43.072-96.736-96-96.736-53.824 0-96 41.536-96 94.56 0 17.664-14.336 32-32 32s-32-14.336-32-32c0-87.424 71.776-158.56 160-158.56s160 72.096 160 160.736c0 47.904-36.32 84.192-71.424 119.296C572.736 532.992 544 561.728 544 587.552l0 54.112C544 659.328 529.664 673.696 512 673.696z" p-id="3680" fill="#8a8a8a"></path></svg>
                    </i>
                  </el-tooltip>
                </el-col>
              </el-row>
              <el-row class="row-item">
                <el-button class="operation-button" v-if="processStatus===0" :disabled="deviceSelected===''" type="primary" @click="handleStartRecord" >开始</el-button>
                <el-button class="operation-button" v-if="processStatus===2" type="danger"  @click="handleStopProcessing" >停止 {{ countDownSecond > 0 ? ": " + countDownSecond : ""}}</el-button>
              </el-row>
              <!-- <el-row class="row-item" v-if="settingForm.develop">
                <el-button-group>
                  <el-button @click="handleStart">rec</el-button>
                  <el-button @click="handleToImage">to_img</el-button>
                  <el-button @click="setPointerLocationOn">pl_on</el-button>
                  <el-button @click="setPointerLocationOff">pl_off</el-button>
                  <el-button @click="handleGetImage">get_imgs</el-button>
                </el-button-group>
              </el-row> -->

              <el-tabs 
                  v-model="latencyTabName" 
                  class="platform-tabs">
                  <el-tab-pane label="记录" name="list">
                    <FileRecord ref="fileRecordRef"/>

                  </el-tab-pane>
                
                  <el-tab-pane label="设置" name="setting">
                    <!-- <el-scrollbar style="height:60vh"> -->
                      <el-row>
                        <el-form :model="settingForm" ref="settingFormRef" :rules="rules" label-position="left" label-width="100px">
                          <el-form-item label="触控阈值" prop="touchScore">
                            <el-input v-model.number="settingForm.touchScore"/>
                          </el-form-item>
                          <el-form-item label="区域阈值" prop="diffScore">
                            <el-input v-model.number="settingForm.diffScore"/>
                          </el-form-item>
                          <el-form-item label="录制时长" prop="timeout">
                            <el-input v-model.number="settingForm.timeout"/>
                          </el-form-item>
                          <el-form-item label="场景时间" prop="sceneStart">
                            <el-input v-model.number="settingForm.sceneStart"/>
                          </el-form-item>
                          <el-form-item label="调试开关">
                            <el-switch v-model="settingForm.develop" />
                          </el-form-item>
                          <el-form-item label="数据清理">
                            <el-button @click="handleClearCache">清理缓存数据</el-button>
                          </el-form-item>
                          <el-form-item label="调式">
                            <el-button @click="handleReload">重载页面</el-button>
                          </el-form-item>
                          <el-form-item label="自动上传">
                            <el-switch v-model="settingForm.autoUpload" />
                          </el-form-item>
                          <el-form-item label="视频地址">
                            <el-input v-model="externalVideoPath"></el-input>
                            <el-button @click="handleLoadExtVideo">加载</el-button>
                          </el-form-item>
                        </el-form>
                      </el-row>
                    <!-- </el-scrollbar> -->
                  </el-tab-pane>
                  <el-tab-pane label="帮助" name="detail" disabled>
                      <HelpPage></HelpPage>
                  </el-tab-pane>
                 
                </el-tabs>
            </el-aside>
            <el-main class="main-content">
              <ImagePreview
                ref="imagePreviewRef"
                :data="imageInfo"
              />
            </el-main>
          </el-container>
        </el-scrollbar>
      </el-tab-pane>
      <el-tab-pane label="场景配置" name="automation">
        <Automation/>
      </el-tab-pane>
      <el-tab-pane label="关于" name="about">
        <AboutPage></AboutPage>
      </el-tab-pane>
    </el-tabs>
   
  </el-container>
    
</template>

<style scoped>

.operation-button {
  width: 100%;
}

.row-item {
  margin-bottom: 7px;
}

.el-form-item {
  margin-bottom: 7px;
}
.aside-content {
  border: solid 1px #e6e6e6;
  padding: 0.5rem;
  border-radius: 4px;
  height: 80vh;
  /* box-shadow: 0 0 6px RGBA(0, 0, 0, 0.2); */
}

.main-content {
  border: solid 1px #e6e6e6;
  padding: 0.5rem;
  border-radius: 4px;
  margin-left: 1rem;
  width: calc(100vw - 320px);
  /* box-shadow: 0 0 6px RGBA(0, 0, 0, 0.2); */
}

.record-list {
  border: solid 1px #e6e6e6;
  padding: 0.5rem;
  border-radius: 4px;
  /* height: 80vh; */
}
.describe {
  opacity: 0.75;
}

.button-icon {
  width: 24px;
  height: 24px;
  margin: 2px;
}
</style>
