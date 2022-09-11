<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { UserFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { 
  ListDevices,
  Start,
  StartRecord,
  StopRecord,
  StopProcessing,
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
  InputSwipe
} from '../../wailsjs/go/app/Api'
import {adb, core} from '../../wailsjs/go/models'
import ImagePreview from '../components/ImagePreview.vue';
import Automation from '../components/Automation.vue';
import {
  EventsOn,
  EventsOff,
  WindowReload,
} from '../../wailsjs/runtime/runtime'

const deviceSelected = ref("")
const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})


const isAuth = ref(false)
const placeholder = "./src/assets/images/placeholder.png"

const rightContentRef = ref()
const startButtonText = ref("开始")
const interval = ref()
const processStatus = ref(0)
const developMode = ref(true)
const countDownSecond = ref(0)
const imageSrc = ref()
const imagePreviewRef = ref()

const topTabName = ref('latency')
// const topTabName = ref('automation')

const tabName = ref('setting')
const deviceInfo = reactive({
  android_version: null,
  cpu_arch: '',
  cpu_core_count: null,
  hardware: '',
  mem_total: 0,
  openGLES_version: '',
  product_model: '',
  width: 0,
  height: 0,
})

const imageInfo = reactive({
  // path: '/Users/jason/Developer/epc/op-latency-mobile/build/bin/op-latency-mobile.app/Contents/MacOS/2022-08-31T14:29:46+08:00/images/0001.png',
  // path: '/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/0001.png',
  path: placeholder,
  width: 0,
  height: 0,
  size: 0
})

const settingForm = reactive({
  diffScore: 20,
  timeout: 3,
  prepareTimeout: 3,
  develop: false,
})

// const threshold = ref(20)

provide('threshold', settingForm.diffScore)

const rules =  {
  diffScore: [
    {
      required: true,
      message: '图片比对阈值',
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
  prepareTimeout: [
    {
      required: true,
      message: '准备时长',
      trigger: 'blur',
    },
    {
      validator: checkGreaterEqualZero,
      trigger: 'blur',
    }
  ]
}

function checkGreaterThanZero (rule: any, value: any, callback: any)  {
  if (value <= 0) {
    callback(new Error('数值必须大于0'))
  } else {
    callback()
  }
}

function checkGreaterEqualZero (rule: any, value: any, callback: any)  {
  if (value < 0) {
    callback(new Error('数值必须大于等于0'))
  } else {
    callback()
  }
}



function getDeviceList (value: any) {
  ListDevices().then(result => {
    if (result != null) {
      data.devices = result
    }
  })
}

const form = reactive({
  name: 'Hello',
  region: '',
  count: '',
  date1: '',
  date2: '',
  delivery: false,
  type: [],
  resource: '',
  desc: '',
  devices: [],
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

const status = ref(false)

// async function handleStartRun() {
//   // await handleGetDisplay()
//   if (form.location) {
//     setPointerLocationOn()
//     // if (!res) {
//     //   return
//     // }
//   }

//   const swipeEvent = adb.SwipeEvent.createFrom(
//     { 
//       sx: form.sx,
//       sy: form.sy,
//       dx: form.dx,
//       dy: form.dy,
//       speed: form.speed
//     }
//   )
//   // const interval = 2
//   console.log(swipeEvent)
//   SetAutoSwipeOn(swipeEvent, form.interval)
//   status.value = 1
// }

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
  // clear first
  // clearCurrentInterval()
  handleResetStatus()
  NProgress.start()
  const result = await setPointerLocationOn()

  await handleGetDisplay()

  Start(deviceSelected.value, settingForm.timeout)

  // 1s 后拖动
  // doSwipe()

  // add stop count down
  processStatus.value = 2
  runUntilCountDown(settingForm.timeout)

  setTimeout(() => {
    handleInputSwipe()
  }, 1200);
}

function handleStart() {
  StartRecord(deviceSelected.value)
}

function handleStopRecord() {
  // clear first
  clearCurrentInterval()

  StopRecord(deviceSelected.value)
  processStatus.value = 0
  SetPointerLocationOff(deviceSelected.value)
}


function handleStopProcessing() {
  StopProcessing()
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

function handleResetStatus() {
  if (NProgress.isStarted()) {
    NProgress.done()
  }
  imagePreviewRef.value.setCalcButtonDisable(true)
  imagePreviewRef.value.setImagePlaceHolder()
  imagePreviewRef.value.setDefaultTime()
}


async function handlePrepare(){
  if (deviceSelected.value === "") {
    ElMessage({
      type: 'error',
      message: '请选择设备'
    })
    return
  }
  handleResetStatus()
  NProgress.start()
  const result = await setPointerLocationOn()
  if (result) {
    if (settingForm.prepareTimeout === 0) {
      processStatus.value = 2
      handleStartRecord()
    } else {
      processStatus.value = 1
      runUntilCountDown(settingForm.prepareTimeout, handleStartRecord)
    }
  }
}

async function setPointerLocationOn():Promise<Boolean> {
  let result = false
  await SetPointerLocationOn(deviceSelected.value).then(res =>{ 
    if (res) {
      ElMessage({
        type: 'error',
        message: '开启指针失败'
      })
      result = false
      
    } else {
      ElMessage({
        type: 'success',
        message: '开启指针成功'
      })
     result = true
    }
  })
  return result
}

function setPointerLocationOff():Boolean {
  SetPointerLocationOff(deviceSelected.value).then(res =>{ 
    if (res) {
      ElMessage({
        type: 'error',
        message: '关闭指针失败'
      })
      return false
      
    } else {
      ElMessage({
        type: 'success',
        message: '关闭指针成功'
      })
      return true
    }
  })
  return false
}

function addEventLister() {
  EventsOn("latency:record_start", ()=>{
    console.log("record_start")
    ElNotification({
      title: '进度提示',
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
      title: '进度提示',
      type: 'success',
      message: "数据预处理中，请稍后...",
      duration: 6000,
    })
    NProgress.done()
  })
  EventsOn("latency:transform_start_error", ()=>{
    ElNotification({
      title: '进度提示',
      type: 'error',
      message: "数据预处理失败，请重试",
      duration: 0,
    })
    NProgress.done()
  })
  EventsOn("latency:record_start_error", ()=>{
    ElNotification({
      title: '进度提示',
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

async function handleGetDisplay() {
  await GetDisplay(deviceSelected.value).then((res: adb.Display) => {
    console.log(res)
    if (res) {
      deviceInfo.width = res.width
      deviceInfo.height = res.height
    } else {
      ElNotification({
        title: '获取数据异常',
        type: 'error',
        message: "获取手机分辨率失败，请手动设置移动坐标",
        duration: 0,
      })
    }
  })
}
async function setAutoOn() {
  await handleGetDisplay()
  const swipeEvent = adb.SwipeEvent.createFrom(
    { 
      sx: deviceInfo.height/2,
      sy: deviceInfo.width/2,
      dx: deviceInfo.height/2 + deviceInfo.height/2/2,
      dy: deviceInfo.width/2,
      speed: 500
    }
  )
  const interval = 2
  console.log(swipeEvent)
  SetAutoSwipeOn(swipeEvent, interval)
}

function setAutoOff() {
  SetAutoSwipeOff()
  
}

function removeEventLister() {
  EventsOff("latency:record_start")
  EventsOff("latency:record_filish")
  EventsOff("latency:transform_start")
  EventsOff("latency:transform_start_error")
  EventsOff("latency:record_start_error")
  EventsOff("latency:transform_filish")
}


function initCheck() {
  console.log("init check")
}

onMounted(()=> {
  initCheck()
  addEventLister()
  
})
function handleClearCache() {
  ClearCacheData()
}

function handleReload() {
  WindowReload();
}

onUnmounted(()=>{
  removeEventLister()
})

function handleGetImage() {
  imagePreviewRef.value.handleGetImage()
}

</script>

<template>
  <el-container>
    <el-tabs type="border-card" v-model="topTabName" >
      <el-tab-pane label="延迟测试" name="latency">
        <el-scrollbar style="height: calc(100vh - 120px);width: calc(100vw - 60px)">
          <el-container>
            <el-aside class="aside-content" width="220px">
              <el-row class="row-item" v-if="isAuth">
                <el-avatar :icon="UserFilled" />
              </el-row>
              <el-row class="row-item">
                <el-select
                    v-model="deviceSelected"
                    @visible-change="getDeviceList"
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
              </el-row>
              <el-row class="row-item">
                <el-button class="operation-button" v-if="processStatus===0" :disabled="deviceSelected===''" type="primary" @click="handleStartRecord" >开始</el-button>
                <!-- <el-button class="operation-button" v-if="processStatus===0" :disabled="deviceSelected===''" type="primary" @click="handlePrepare" >准备</el-button> -->
                <!-- <el-button class="operation-button" v-if="processStatus===1" type="success" @click="handleStartRecord" >开始 {{ countDownSecond > 0 ? ": " + countDownSecond : ""}}</el-button> -->
                <el-button class="operation-button" v-if="processStatus===2" type="danger"  @click="handleStopProcessing" >停止 {{ countDownSecond > 0 ? ": " + countDownSecond : ""}}</el-button>
              </el-row>
              <el-row class="row-item" v-if="settingForm.develop">
                <el-button-group>
                  <el-button @click="handleStart">rec</el-button>
                  <el-button @click="handleStopProcessing">stop</el-button>
                  <el-button @click="handleToImage">to_img</el-button>
                  <el-button @click="setPointerLocationOn">pl_on</el-button>
                  <el-button @click="setPointerLocationOff">pl_off</el-button>
                  <el-button @click="setAutoOn">auto_on</el-button>
                  <el-button @click="setAutoOff">auto_off</el-button>
                  <el-button @click="handleGetImage">get_imgs</el-button>
                </el-button-group>
              </el-row>

              <el-tabs 
                  v-model="tabName" 
                  class="platform-tabs">
                
                  <el-tab-pane label="设置" name="setting">
                    <el-scrollbar style="height:60vh">
                      <el-row>
                        <el-form :model="settingForm" ref="settingFormRef" :rules="rules">
                          <el-form-item label="图片比对阈值" prop="diffScore">
                            <el-input v-model.number="settingForm.diffScore"/>
                            <!-- <el-input-numbe v-model="settingForm.diffScore" :min="1" :max="100" ></el-input-numbe> -->
                          </el-form-item>
                          <el-form-item label="录制时长(秒)" prop="timeout">
                            <el-input v-model.number="settingForm.timeout"/>
                          </el-form-item>
                          <!-- <el-form-item label="准备时长(秒)" prop="prepareTimeout">
                            <el-input v-model.number="settingForm.prepareTimeout"/>
                          </el-form-item> -->
                          <el-form-item label="调试开关">
                            <el-switch v-model="settingForm.develop" />
                          </el-form-item>
                          <el-form-item label="数据清理">
                            <el-button @click="handleClearCache">清理缓存数据</el-button>
                          </el-form-item>
                          <el-form-item label="其他">
                            <el-button @click="handleReload">reload</el-button>
                          </el-form-item>
                        </el-form>
                      </el-row>
                    </el-scrollbar>
                  </el-tab-pane>
                  <el-tab-pane label="帮助" name="detail">
                    <el-scrollbar style="height:60vh">
                      <div>
                        <!-- <el-row class="info-list">
                          <el-col :span="12" class="info-line">
                            <span class="info-key">名称</span>
                          </el-col >
                          <el-col :span="12" class="info-line">
                            <span class="info-value">数值</span>
                          </el-col>
                          <el-row class="info-list">
                          <el-col :span="12" class="info-line">
                            <span class="info-key">Version</span>
                          </el-col >
                          <el-col :span="12" class="info-line">
                            <span class="info-value">{{ deviceInfo.android_version }}</span>
                          </el-col>
                        </el-row>
                        </el-row> -->
                        
                      </div>
                    </el-scrollbar>
                  </el-tab-pane>
                  <el-tab-pane label="关于" name="about">
                    <span>www.vrviu.com</span>
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
      <el-tab-pane label="滑动工具" name="automation" disabled>
        <el-scrollbar style="height:calc(100vh - 120px);width: calc(100vw - 60px)">
          <Automation/>
        </el-scrollbar>
      </el-tab-pane>
    </el-tabs>
   
  </el-container>
    
</template>

<style scoped>
/* .v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}


.v-enter-from,
.v-leave-to {
  opacity: 0;
} */

.operation-button {
  width: 100%;
}

.row-item {
  margin-bottom: 1rem;
}

.aside-content {
  border: solid 1px #e6e6e6;
  padding: 0.5rem;
  border-radius: 4px;
  height: 86vh;
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
</style>
