<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { UserFilled, VideoPlay } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import ScreenPreview from '../components/ScreenPreview.vue';

import { isWailsRun } from '@/utils/utils'

import { 
  ListDevices,
  Start,
  StartRecord,
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
  IsPointerLocationOn,
  InputTap,
} from '../../wailsjs/go/app/Api'
import {adb, core} from '../../wailsjs/go/models'
import {
  EventsOn,
  EventsOff,
  WindowReload,
} from '../../wailsjs/runtime/runtime'
import { stat } from 'fs'

const sceneSelected = ref("")
const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})

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

const latencyForm = reactive({
  serial: '',
  auto: true,
  scene: ''
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
  // path: '/Users/jason/Developer/epc/op-latency-mobile/tmp001.png',
  path: '',
  width: 0,
  height: 0,
  size: 0
})

const result = reactive({
  latency: 0,
  responseIndex: 0,
  responseTime: 0,
  imageCount: 0,
  inputTime: 0,
})
const imageZoom = ref(1)

const cropInfo:CropArea = reactive({
  top: 50,
  left: 50,
  width: 90,
  height: 90,
})

const userAction = reactive({
  type: 'click',
  x: 0,
  y: 0,
  tx: 0,
  ty:0,
  speed: 0,
})

const realCropInfo:Ref<CropArea> = computed(()=> {
  return {
    top: cropInfo.top * imageZoom.value,
    left: cropInfo.left * imageZoom.value,
    width: cropInfo.width * imageZoom.value,
    height: cropInfo.height * imageZoom.value,
  }
})

const realUserAction = computed(() => {
  return {
    type: 'click',
    x: userAction.x * imageZoom.value,
    y: userAction.y * imageZoom.value,
    tx: userAction.tx * imageZoom.value,
    ty: userAction.ty * imageZoom.value,
    speed: userAction.speed,
  }
})


const imagePageInfo:ImagePage = reactive({
  size: 1,
  total: 0,
  currentPage: 1,
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
    if (d.Serial == latencyForm.serial) {
      return d.State
    }
  }
}


// function handleLoadExtVideo() {
//   handleResetStatus()
//   NProgress.start()
//   StartWithVideo(externalVideoPath.value)
// }


/**
 * 发送操作事件
 */
 function handleInput() {
  if (userAction.type === 'swipe') {
    const swipeEvent = adb.SwipeEvent.createFrom({ 
      sx: userAction.x,
      sy: userAction.y,
      dx: userAction.tx,
      dy: userAction.ty,
      speed: userAction.speed
    })  
    InputSwipe(latencyForm.serial, swipeEvent)
  } else {
    const tapEvent = adb.TapEvent.createFrom({
      x: userAction.x,
      y: userAction.y
    })
    InputTap(latencyForm.serial, tapEvent).then().catch(err => { console.log(err)})
  }
}


const handleCropChange = (res: CropInfo)=> {
  cropInfo.left = res.left
  cropInfo.top = res.top
  cropInfo.width = res.width
  cropInfo.height = res.height
}

const handlePageChange = (val: number) => {
  imagePageInfo.currentPage = val
  handleLoadImage(imagePageInfo.currentPage -1)
}

const handleOpenFolder = (val: number) => {
  // OpenImageInExplorer(val).then().catch(err => console.log(err))
}

const handleUserAction = (action: any) => {
  console.log(action)

  if (action.type === 'click') {
    userAction.type = action.type
    userAction.x = action.x
    userAction.y = action.y
  } else {
    userAction.type = action.type
    userAction.x = action.x
    userAction.y = action.y
    userAction.tx = action.tx
    userAction.ty = action.ty
    userAction.speed = action.speed
  }
}

const handleLoadImage = (val: number) => {}
const handleCalcWithCurrent = () => {}

// function runUntilCountDown(second: number, callback?: Function){
//   countDownSecond.value = second
//   function countDown() {
//     // 启动需要时间，提前1s启动
//     if (countDownSecond.value  == 1) {
//       if (callback) { callback() }
//     }
//     if (countDownSecond.value  > 0) {
//       countDownSecond.value  --
//     }

//   }

//   clearCurrentInterval()
//   interval.value = setInterval(countDown, 1000)
// }

/**
 * 配置默认状态
 */
// function handleResetStatus() {
//   if (NProgress.isStarted()) {
//     NProgress.done()
//   }
//   imagePreviewRef.value.setCalcButtonDisable(true)
//   imagePreviewRef.value.setImagePlaceHolder()
//   imagePreviewRef.value.setDefaultTime()
// }


// function getFirstImage(){
//   GetFirstImageInfo().then((res: core.ImageInfo) => {
//     imageInfo.path = res.path
//     imageInfo.width = res.width
//     imageInfo.height = res.height
//     imagePreviewRef.value.loadNewImage(res)
//     imagePreviewRef.value.enableCalcButton()
//   })
// }


// function handleClearCache() {
//   ClearCacheData()
// }

// function handleReload() {
//   WindowReload();
// }

// function handleStopProcessing() {

// }

// function handleGetImage() {
//   imagePreviewRef.value.handleGetImage()
// }


watch(imageInfo, (val: any) => {
  const pImgSize = imagePreviewRef.value.getPreviewImgSize()
  imageZoom.value = imageInfo.width / pImgSize.width
})


onMounted(()=> {
  // 如果是在 wails 运行环境则运行环境检查及事件监听
  if (isWailsRun()) {
    // initCheck()
    // addEventLister()
    // fileRecordRef.value.handleLoadCacheFiles()
  }
})

onUnmounted(()=>{
  if (isWailsRun()) {
    // removeEventLister()
  }
})


</script>

<template>
    <el-scrollbar style="height: calc(100vh - 100px);width: calc(100vw - 60px)">
        <el-container>
        <el-aside class="aside-content" width="240px">
            <el-row class="row-item">
              <el-form :model="latencyForm">
                <el-form-item label="设备">
                  <el-col :span="20">
                  <el-select
                    v-model="latencyForm.serial"
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
                </el-form-item>
                <el-form-item label="自动">
                  <el-switch v-model="latencyForm.auto" />
                </el-form-item>
                <el-form-item v-if="latencyForm.auto===true" label="场景">
                  <el-col :span="20">
                    <el-select
                      v-model="latencyForm.scene"
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
                </el-form-item>
              </el-form>
            </el-row>
            <el-row class="row-item">
            </el-row>
            <!-- <el-button @click="handleReload">重载页面</el-button> -->
        </el-aside>
        <el-main class="main-content">
          <div>
              <ScreenPreview
                ref="imagePreviewRef"
                :imageInfo="imageInfo"
                :cropInfo="cropInfo"
                :pageInfo="imagePageInfo"
                @crop-change="handleCropChange"
                @page-change="handlePageChange"
                @open-folder="handleOpenFolder"
                @user-action="handleUserAction"
                />
            </div>

            <el-row justify="center" class="result-row">
            </el-row>
            <el-row justify="center" class="result-row">
              <el-col :span="4" class="info-line">
                <span>动作</span>
              </el-col>
              <el-col :span="8" class="info-line">
                <span v-if="userAction.type==='click'">点击 x: {{ realUserAction.x}} y: {{ realUserAction.y }}</span>
                <span v-if="userAction.type==='swipe'">滑动 x: {{ realUserAction.x}} y:{{ realUserAction.y}} tx: {{ realUserAction.tx }} ty:{{ realUserAction.ty }}  speed: {{ realUserAction.speed }}</span>
              </el-col>
            </el-row>
            <el-row justify="center" class="result-row">
              <el-col :span="4" class="info-line">
                <span>观察区域</span>
              </el-col>
              <el-col :span="8" class="info-line">
                left: {{ realCropInfo.left }} top: {{ realCropInfo.top }}  width: {{ realCropInfo.width }}  height: {{ realCropInfo.height }}
              </el-col>
            </el-row>
            <el-row justify="center" class="button-row">
              <el-button type="success" @click="handleInput" :icon="VideoPlay">
                执行动作
              </el-button>
              <el-button @click="handleCalcWithCurrent">
                <i class="el-icon button-icon">
                  <svg t="1666605626827" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4819" width="200" height="200"><path d="M665.6 332.8 665.6 128C665.6 113.86151 654.13849 102.4 640 102.4 625.86151 102.4 614.4 113.86151 614.4 128L614.4 332.8C614.4 346.93849 625.86151 358.4 640 358.4 654.13849 358.4 665.6 346.93849 665.6 332.8L665.6 332.8ZM640 51.2 819.2 51.2 793.6 25.6 793.6 384.133545C793.6 426.406699 759.102946 460.8 716.727898 460.8L281.672102 460.8C239.236715 460.8 204.8 426.413438 204.8 384.133545L204.8 25.6C204.8 11.46151 193.33849 0 179.2 0 165.06151 0 153.6 11.46151 153.6 25.6L153.6 384.133545C153.6 454.707134 210.976425 512 281.672102 512L716.727898 512C787.345461 512 844.8 454.718257 844.8 384.133545L844.8 25.6 844.8 0 819.2 0 640 0C625.86151 0 614.4 11.46151 614.4 25.6 614.4 39.73849 625.86151 51.2 640 51.2L640 51.2Z" p-id="4820" fill="#8a8a8a"></path><path d="M844.8 972.8 128.081132 972.8C85.544157 972.8 51.2 938.575806 51.2 896.163853L51.2 100.711064 51.2 25.6 25.6 51.2 102.4 51.2 896.233363 51.2C938.580175 51.2 972.8 85.414085 972.8 127.868001L972.8 998.4C972.8 1012.53849 984.26151 1024 998.4 1024 1012.53849 1024 1024 1012.53849 1024 998.4L1024 127.868001C1024 57.135182 966.85523 0 896.233363 0L102.4 0 25.6 0 0 0 0 25.6 0 100.711064 0 896.163853C0 966.892966 57.307204 1024 128.081132 1024L844.8 1024C858.93849 1024 870.4 1012.53849 870.4 998.4 870.4 984.26151 858.93849 972.8 844.8 972.8L844.8 972.8Z" p-id="4821" fill="#8a8a8a"></path></svg>
                </i>
                保存 </el-button>
            </el-row>
        </el-main>
        </el-container>
    </el-scrollbar>
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

.info-line {
  border: solid 1px #e6e6e6;
  opacity: 0.6;
  font-size: 12px;
  line-height: 18px;
  display:table-cell;
  vertical-align:middle;
  padding: 3px;
}

.button-row {
  margin: 7px 0;
}
</style>
