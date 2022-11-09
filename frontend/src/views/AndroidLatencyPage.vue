<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

// import ImagePreview from '../components/ImagePreview.vue';
import FileRecord from '../components/FileRecord.vue';
// import Automation from '../components/Automation.vue';
// import AboutPage from '../components/AboutPage.vue';
// import HelpPage from '../components/HelpPage.vue';
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
  ClearMobleCache,
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
  ListScens,
  InputTap,
} from '../../wailsjs/go/app/Api'
import {adb, core, app} from '../../wailsjs/go/models'
import {
  EventsOn,
  EventsOff,
  WindowReload,
} from '../../wailsjs/runtime/runtime'
import { stat } from 'fs'


const deviceSelected = ref("")
const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})

const latencyTabName = ref('list')

const fileRecordRef = ref()
const loadExtVideoVisible = ref(false)
const extVideoPath = ref('')

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
  auto: true,
  scene: {
    name: '',
    action: {
      type: 'swipe',
      x: 0,
      y: 0,
      tx: 0,
      ty: 0,
      speed: 500
    },
    crop_coordinate: {
      left: 100,
      top: 100,
      width: 50,
      height:50
    },
    crop_touch_coordinate: {
      left: 0,
      top: 0,
      width: 100,
      height:35
    },
  },
  device: {
    serial: '',
    state: 0,
    deviceName: '',
  }
})

const interval = ref()
const processStatus = ref(0)
const countDownSecond = ref(0)
const imagePreviewRef = ref()
const externalVideoPath = ref('')
const calcButtonDisable = ref(true)
const imageZoom = ref(1)
const deviceInfo = reactive({
  width: 1080,
  height: 1920,
})

const imageInfo = reactive({
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

const cropInfo:CropArea = reactive({
  top: 100,
  left: 100,
  width: 50,
  height: 50,
})

const cropTouchInfo:CropArea = reactive({
  top: 100,
  left: 100,
  width: 50,
  height: 50,
})

const userScenes: {scens: Array<app.UserScene>} = reactive({
  scens: [],
})

// @ts-ignore:  default value
// const userScene:app.UserScene = reactive({
//   name: '',
//   device: {},
//   crop_coordinate: {
//     top: 100,
//     left: 100,
//     width: 50,
//     height: 50,
//   },
//   action: {},
// })

const selectedScene = ref<app.UserScene>()

const imagePageInfo:ImagePage = reactive({
  size: 1,
  total: 0,
  currentPage: 1,
})

const settingForm = reactive({
  touchScore: 1,
  diffScore: 20,
  timeout: 1,
  sceneStart: 1500,
  prepareTimeout: 3,
  develop: false,
  autoUpload: false,
  // pointLocation: true,
})

provide('threshold', settingForm.diffScore)

const imagePreviewInfo = reactive({
  width: 0,
  height: 0,
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
 * 设备选中时检测设备状态
 */
const handleDeviceChange = () => {
  const state = latencyForm.device.state
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
}

/**
 * 获取选中设备状态
 */
// function getSelectDeviceState(){
//   for(let d of data.devices) {
//     if (d.serial == latencyForm.serial) {
//       return d.state
//     }
//   }
// }

const handleGetScenes = () => {
  ListScens().then((res: Array<app.UserScene>) => {
    console.log(res)
    userScenes.scens = res
  }).catch(err => {
    console.log(err)
  })
}



const handleSceneChange = () => {
  console.log('handleSceneChange')
  if (selectedScene.value != undefined || selectedScene.value != null) {
    latencyForm.scene.name = selectedScene.value.name
    latencyForm.scene.crop_coordinate = selectedScene.value.crop_coordinate
    latencyForm.scene.crop_touch_coordinate = selectedScene.value.crop_touch_coordinate
    latencyForm.scene.action = selectedScene.value.action
    
      // userScene.crop_coordinate = val.crop_coordinate
  // userScene.action = val.action

  // imagePreviewRef.value.selectBoxInit()
  // cropInfo.left = val.crop_coordinate.left
  // cropInfo.top = val.crop_coordinate.top
  // cropInfo.width = val.crop_coordinate.width
  // cropInfo.height = val.crop_coordinate.height
  }
  // console.log(latencyForm.scene)
  // userScene.name = val.name
  // userScene.device = val.device
  // userScene.crop_coordinate = val.crop_coordinate
  // userScene.action = val.action

  // imagePreviewRef.value.selectBoxInit()
  // cropInfo.left = val.crop_coordinate.left
  // cropInfo.top = val.crop_coordinate.top
  // cropInfo.width = val.crop_coordinate.width
  // cropInfo.height = val.crop_coordinate.height

  // console.log(userScene)
}

// const deviceStateCheck = () => {
  
// }

function handleLoadExtVideo() {
  // handleResetStatus()
  // NProgress.start()
  // StartWithVideo(externalVideoPath.value)
  loadExtVideoVisible.value = true
}

function handleStartWithExtVideo() {
  handleResetStatus()
  NProgress.start()
  StartWithVideo(extVideoPath.value)
  loadExtVideoVisible.value = false
}

/**
 * 发送拖动事件
 */
// function handleInputSwipe() {
//   const swipeEvent = adb.SwipeEvent.createFrom(
//     { 
//       sx: Math.trunc(deviceInfo.height/2),
//       sy: Math.trunc(deviceInfo.width/2),
//       dx: Math.trunc(deviceInfo.height/2) + Math.trunc(deviceInfo.height/2/2),
//       dy: Math.trunc( deviceInfo.width/2),
//       speed: form.speed
//     }
//   )
//   // const interval = 2
//   console.log(swipeEvent)
//   InputSwipe(latencyForm.device.serial, swipeEvent)
// }

/**
 * 发送操作事件
 */
 function handleAutoInput() {
  console.log("handleAutoInput")
  if (selectedScene.value === undefined) {
    ElMessage({
      type: 'warning',
      message: '场景信息为空',
    })
    return
  }
  if (selectedScene.value.action.type === 'swipe') {
    const swipeEvent = adb.SwipeEvent.createFrom({ 
      sx: selectedScene.value.action.x,
      sy: selectedScene.value.action.y,
      dx: selectedScene.value.action.tx,
      dy: selectedScene.value.action.ty,
      speed: selectedScene.value.action.speed
    })  
    InputSwipe(latencyForm.device.serial, swipeEvent).then().catch(err => { console.log(err)})
  } else {
    const tapEvent = adb.TapEvent.createFrom({
      x: selectedScene.value.action.x,
      y: selectedScene.value.action.y
    })
    InputTap(latencyForm.device.serial, tapEvent).then().catch(err => { console.log(err)})
  }

  // if (latencyForm.scene.action.type === 'swipe') {
  //   const swipeEvent = adb.SwipeEvent.createFrom({ 
  //     sx: latencyForm.scene.action.x,
  //     sy: latencyForm.scene.action.y,
  //     dx: latencyForm.scene.action.tx,
  //     dy: latencyForm.scene.action.ty,
  //     speed: latencyForm.scene.action.speed
  //   })  
  //   InputSwipe(latencyForm.device.serial, swipeEvent).then().catch(err => { console.log(err)})

  // } else if (latencyForm.scene.action.type === 'click') {
  //   const tapEvent = adb.TapEvent.createFrom({
  //     x: latencyForm.scene.action.x,
  //     y: latencyForm.scene.action.y
  //   })
  //   InputTap(latencyForm.device.serial, tapEvent).then().catch(err => { console.log(err)})
  // }
}


/**
 * 启动
 */
async function handleStart() {
  if (latencyForm.auto === true && Object.keys(latencyForm.scene).length === 0) {
    ElMessage({
      type: 'warning',
      message: '勾选自动操作后需选择场景',
    })
    return
  }
  //   setTimeout(() => {
  //     handleAutoInput()
  //     // handleInputSwipe()
  //   }, 1500);
  // }
  
  // 重置状态，开启进度条
  handleResetStatus()
  NProgress.start()

  // 查看指针开启状态

  const pStatus = await isPointerLocationOn()
  if (!pStatus) {
    await setPointerLocationOn()
  }
  
  // 获取屏幕分辨率
  // const status = await handleGetPhysicalSize()
  // if (!status) {
  //   await handleGetDisplay()
  // }
  // if (latencyForm.auto === true) {
  const action = app.UserAction.createFrom({
    auto: latencyForm.auto,
    type: latencyForm.scene.action.type,
    x: latencyForm.scene.action.x,
    y: latencyForm.scene.action.y,
    tx: latencyForm.scene.action.tx,
    ty: latencyForm.scene.action.ty,
    speed: latencyForm.scene.action.speed,
  })
  // }
  Start(latencyForm.device.serial, settingForm.timeout, action)

  // 1s 后拖动
  // doSwipe()

  // add stop count down
  processStatus.value = 2
  runUntilCountDown(settingForm.timeout)

  // if (latencyForm.auto === true) {
  //   setTimeout(() => {
  //     handleAutoInput()
  //     // handleInputSwipe()
  //   }, 1500);
  // }
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

const handleCropTouchChange = (res: CropInfo)=> {
  cropTouchInfo.left = res.left
  cropTouchInfo.top = res.top
  cropTouchInfo.width = res.width
  cropTouchInfo.height = res.height
}

const handleGetPreviewImage = () => {
  const pImgSize = imagePreviewRef.value.getPreviewImgSize()
  imagePreviewInfo.width = pImgSize.width
  imagePreviewInfo.height = pImgSize.height

  imageZoom.value = imageInfo.width / pImgSize.width
}

const handleScaleChange = () => {
  // const pImgSize = imagePreviewRef.value.getPreviewImgSize()
  // imageZoom.value = imageInfo.width / pImgSize.width

  handleGetPreviewImage()

  imagePreviewRef.value.updateSelectBoxStyle()
  imagePreviewRef.value.updateSelectBoxTouchStyle()
}


const handleLoadImage = (val: number) => {}

const handleCalc = () => {
  if (selectedScene.value === undefined) {
    ElMessage({
      type: 'warning',
      message: '场景信息为空',
    })
    return
  }

 

  // const pImgSize = imagePreviewRef.value.getPreviewImgSize()
  handleGetPreviewImage()
  const rectinfo = core.ImageRectInfo.createFrom({
    x: selectedScene.value.crop_coordinate.left,
    y: selectedScene.value.crop_coordinate.top,
    w: selectedScene.value.crop_coordinate.width + selectedScene.value.crop_coordinate.left,
    h: selectedScene.value.crop_coordinate.height + selectedScene.value.crop_coordinate.top,
    preview_width: imagePreviewInfo.width,
    preview_height: imagePreviewInfo.height,
    source_width: imageInfo.width,
    source_height: imageInfo.height,
  })

  const rectTouchinfo = core.ImageRectInfo.createFrom({
    x: selectedScene.value.crop_touch_coordinate.left,
    y: selectedScene.value.crop_touch_coordinate.top,
    w: selectedScene.value.crop_touch_coordinate.width + selectedScene.value.crop_touch_coordinate.left,
    h: selectedScene.value.crop_touch_coordinate.height + selectedScene.value.crop_touch_coordinate.top,
  })
  // const threshold = 20
  
  const threshold = app.Threshold.createFrom({
    pointer_threshold: settingForm.touchScore,
    black_white_threshold: 60,
    scene_threshold: settingForm.diffScore,
  })
  StartAnalyse(rectinfo, rectTouchinfo, threshold)
  handleResetStatus()
  // NProgress.start()
  // delayTimes.value = 0 
  // calcButtonDisable.value = true
}
const handleCalcWithCurrent = () => {}


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

  result.latency = 0
  // imagePreviewRef.value.setCalcButtonDisable(true)
  // imagePreviewRef.value.setImagePlaceHolder()
  // imagePreviewRef.value.setDefaultTime()
}


/**
 * 开启指针显示
 */
async function setPointerLocationOn():Promise<Boolean> {
  let result = false
  await SetPointerLocationOn(latencyForm.device.serial).then(res =>{ 
    ElMessage({
      type: 'success',
      message: '开启指针成功'
    })
     result = true
    }
  ).catch(err => {
    ElMessage({
      type: 'warning',
      message: '开启指针失败, 请确定已手工开启'
    })
    result = false
  })
  return result
}

const isPointerLocationOn = async() => {
  let status = false 
  IsPointerLocationOn(latencyForm.device.serial).then((res: boolean) => {
    status = res
  }).catch(err => {
    console.log(err)
  })
  return status
}

function setPointerLocationOff():Boolean {
  SetPointerLocationOff(latencyForm.device.serial).then(res =>{ 
      ElMessage({
        type: 'success',
        message: '关闭指针成功'
      })
      return true
  }).catch(err => {
    ElMessage({
      type: 'warning',
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
    })
    NProgress.done()
  })
  EventsOn("latency:record_start_error", ()=>{
    ElNotification({
      title: '进度提示: 2/3',
      type: 'error',
      message: "录制失败，请重试",
    })
    NProgress.done()
  })
  EventsOn("latency:transform_filish", async()=>{
    ElNotification({
      title: '进度提示: 2/3',
      type: 'success',
      message: "数据预处理完成，加载首帧画面",
    })
    await getFirstImage()

    NProgress.done()
    updateCropInfo()
    fileRecordRef.value.handleLoadCacheFiles()
  })
  EventsOn("latency:analyse_start", ()=>{
    ElNotification({
      title: '进度提示',
      type: 'info',
      message: "数据分析中， 请稍后...",
    })
  })

  EventsOn("latency:record_file_exists", () => {
    if (latencyForm.auto === true) {
      ElNotification({
        title: '操作提醒',
        message: "开始自动化操作",
      })
      // handleAutoInput()
    }
  })
    

  EventsOn("latency:analyse_filish", (res: number)=>{
    if (res) {
      result.latency =  Math.floor(res * 100)/100
      ElNotification({
        title: '进度提示: 3/3',
        type: 'success',
        message: "数据处理完成",
      })
      NProgress.done()
      calcButtonDisable.value = false

      if (result.latency <= 50 || result.latency >= 1000 ) {
        ElNotification({
          title: '数值异常',
          type: 'error',
          message: "当前数值不在串流延迟正常范围内，建议重试",
        })
      }
    } else {
      NProgress.done()
      calcButtonDisable.value = false

      ElNotification({
        title: '进度提示: 3/3',
        type: 'error',
        message: "数据分析异常，请确认是否在指定业务场景下操作，建议重试",
      })
    }
  })
}

const updateCropInfo = () => {
  if (selectedScene.value === undefined && latencyForm.auto === true) {
    ElMessage({
      type: 'warning',
      message: '场景信息为空',
    })
  }

  const pImgSize = imagePreviewRef.value.getPreviewImgSize()
  imageZoom.value = pImgSize.width / imageInfo.width 
  // console.log(imageZoom.value)
  
  cropInfo.left = latencyForm.scene.crop_coordinate.left * imageZoom.value,
  cropInfo.top = latencyForm.scene.crop_coordinate.top * imageZoom.value,
  cropInfo.width = latencyForm.scene.crop_coordinate.width * imageZoom.value,
  cropInfo.height = latencyForm.scene.crop_coordinate.height * imageZoom.value,

  cropTouchInfo.left = latencyForm.scene.crop_touch_coordinate.left * imageZoom.value,
  cropTouchInfo.top = latencyForm.scene.crop_touch_coordinate.top * imageZoom.value,
  cropTouchInfo.width = latencyForm.scene.crop_touch_coordinate.width * imageZoom.value,
  cropTouchInfo.height = latencyForm.scene.crop_touch_coordinate.height * imageZoom.value,

  // console.log(cropInfo)
  imagePreviewRef.value.updateSelectBoxStyle()
  imagePreviewRef.value.switchSelectBoxShow(true)
  imagePreviewRef.value.updateSelectBoxTouchStyle()
  imagePreviewRef.value.switchSelectBoxTouchShow(true)
}

const isVerticalScreen = () => {
  return imageInfo.width < imageInfo.height
}

async function getFirstImage(){
  await GetFirstImageInfo().then((res: core.ImageInfo) => {
    console.log('getFirstImage')
    imageInfo.path = res.path
    imageInfo.width = res.width
    imageInfo.height = res.height

    if (isVerticalScreen()) {
      imagePreviewRef.value.setScalePercent(50)
    } else {
      imagePreviewRef.value.setScalePercent(100)
    }

    calcButtonDisable.value = false
    
  }).catch(err => {
    console.log(err)
  })
}

async function handleGetPhysicalSize() {
  let status = false
  await GetPhysicalSize(latencyForm.device.serial).then((res: adb.Display) => {
      deviceInfo.width = res.width
      deviceInfo.height = res.height
      status = true
  }).catch(err => {
    console.log(err)
  })
  return status
}


async function handleGetDisplay() {
  let status = false
  await GetDisplay(latencyForm.device.serial).then((res: adb.Display) => {
      deviceInfo.width = res.width
      deviceInfo.height = res.height
      status = true
  }).catch(err => {
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
  EventsOff("latency:analyse_start")
  EventsOff("latency:record_file_exists")
  EventsOff("latency:analyse_filish")
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
    })
  })
}

function handleClearCache() {
  ClearMobleCache()
}

function handleReload() {
  WindowReload();
}

function handleStopProcessing() {

}

function handleGetImage() {
  imagePreviewRef.value.handleGetImage()
}


onMounted(()=> {
  // 如果是在 wails 运行环境则运行环境检查及事件监听
  if (isWailsRun()) {
    initCheck()
    addEventLister()
    fileRecordRef.value.handleLoadCacheFiles()
  }
})

onUnmounted(()=>{
  if (isWailsRun()) {
    removeEventLister()
  }
})


</script>

<template>
  <el-container class="panel-container">
    <el-aside class="aside-content" width="240px">
        <el-row class="row-item">
          <el-form :model="latencyForm">
            <el-form-item label="设备">
              <el-col :span="20">
              <el-select
                v-model="latencyForm.device"
                @focus="getDeviceList"
                @change="handleDeviceChange"
                filterable
                placeholder="请选择设备"
                style="width:100%">
                <el-option
                  v-for="item in data.devices"
                  :key="item.serial"
                  :label="item.serial"
                  :value="item"
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
                  v-model="selectedScene"
                  filterable
                  placeholder="请选择场景"
                  @focus="handleGetScenes"
                  @change="handleSceneChange"
                  style="width:100%">
                  <el-option
                    v-for="item in userScenes.scens"
                    :key="item.key"
                    :label="item.name"
                    :value="item"
                  >
                  </el-option>
                </el-select>
              </el-col>
              <el-col :span="4">
                <el-tooltip
                  class="device-question"
                  effect="dark"
                  content="选择自动运行场景, 如列表为空，请在场景配置页签配置"
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
        <el-button class="operation-button" v-if="processStatus===0" :disabled="latencyForm.device.serial===''" type="primary" @click="handleStart" >开始</el-button>
        <el-button class="operation-button" v-if="processStatus===2" type="danger"  @click="handleStopProcessing" >停止 {{ countDownSecond > 0 ? ": " + countDownSecond : ""}}</el-button>
        </el-row>
        <el-tabs 
            v-model="latencyTabName" 
            class="platform-tabs">
            <el-tab-pane label="记录" name="list" class="menu-tab">
              <FileRecord 
               ref="fileRecordRef"
               @load-extVideo="handleLoadExtVideo"
               />
            </el-tab-pane>
        
            <el-tab-pane label="设置" name="setting">
              <el-row>
                <el-form :model="settingForm" ref="settingFormRef" :rules="rules" label-position="left" label-width="100px">
                  <el-form-item label="触控阈值" prop="touchScore">
                    <el-input v-model="settingForm.touchScore"/>
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
                  <!-- <el-form-item label="指针位置">
                    <el-switch v-model="settingForm.pointLocation" />
                  </el-form-item> -->
                  <el-form-item label="自动上传">
                    <el-switch v-model="settingForm.autoUpload" />
                  </el-form-item>
                  <el-form-item label="调式">
                    <el-button @click="handleReload">重载页面</el-button>
                  </el-form-item>
                </el-form>
              </el-row>
            </el-tab-pane>
            <!-- <el-tab-pane label="帮助" name="detail" disabled>
                <HelpPage></HelpPage>
            </el-tab-pane> -->
        </el-tabs>
    </el-aside>
    <el-main class="main-content">
      <el-row justify="center" class="result-row">
        <el-col :span="8" class="info-line">
          <span>动作: </span>
          <span v-if="latencyForm.scene.action.type==='click'">点击 x: {{ latencyForm.scene.action.x}} y: {{ latencyForm.scene.action.y }}</span>
          <span v-if="latencyForm.scene.action.type==='swipe'">滑动 x: {{ latencyForm.scene.action.x}} y:{{ latencyForm.scene.action.y}} tx: {{ latencyForm.scene.action.tx }} ty:{{ latencyForm.scene.action.ty }}  s: {{ latencyForm.scene.action.speed }}</span>
        </el-col>
        <el-col :span="8" class="info-line">
          <span class="touch-area">触控区域:</span>
          x: {{ latencyForm.scene.crop_touch_coordinate.left }} y: {{ latencyForm.scene.crop_touch_coordinate.top }}  w: {{ latencyForm.scene.crop_touch_coordinate.width }}  h: {{ latencyForm.scene.crop_touch_coordinate.height }}
        </el-col>
        <el-col :span="8" class="info-line">
          <span class="watch-area">观察区域:</span>
          x: {{ latencyForm.scene.crop_coordinate.left }} y: {{ latencyForm.scene.crop_coordinate.top }}  w: {{ latencyForm.scene.crop_coordinate.width }}  h: {{ latencyForm.scene.crop_coordinate.height }}
        </el-col>
      </el-row>
      <el-scrollbar style="height: calc(100vh - 130px)">
        <div>
          <ScreenPreview
            ref="imagePreviewRef"
            :imageInfo="imageInfo"
            :cropInfo="cropInfo"
            :cropTouchInfo="cropTouchInfo"
            :pageInfo="imagePageInfo"
            @crop-change="handleCropChange"
            @page-change="handlePageChange"
            @open-folder="handleOpenFolder"
            @crop-touch-change="handleCropTouchChange"
            @scale-change="handleScaleChange"
            />
        </div>
        <el-row justify="center" class="button-row">
          <el-button type="success" :disabled="calcButtonDisable" @click="handleCalc">
            <i class="el-icon button-icon">
              <svg t="1666320784905" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5742" width="200" height="200"><path d="M928 1024H96a96 96 0 0 1-96-96V96a96 96 0 0 1 96-96h832a96 96 0 0 1 96 96v832a96 96 0 0 1-96 96zM896 160a32 32 0 0 0-32-32H160a32 32 0 0 0-32 32v160h768V160z m0 288H128v416a32 32 0 0 0 32 32h704a32 32 0 0 0 32-32V448z m-256 64h128v320h-128V512z m-192 192h128v128h-128v-128z m0-192h128v128h-128v-128z m-192 192h128v128H256v-128z m0-192h128v128H256v-128z" p-id="5743" fill="#8a8a8a"></path></svg>
            </i>
            计算延迟
          </el-button>
          <el-button :disabled="calcButtonDisable"  @click="handleCalcWithCurrent">
            <i class="el-icon button-icon">
              <svg t="1666320784905" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5742" width="200" height="200"><path d="M928 1024H96a96 96 0 0 1-96-96V96a96 96 0 0 1 96-96h832a96 96 0 0 1 96 96v832a96 96 0 0 1-96 96zM896 160a32 32 0 0 0-32-32H160a32 32 0 0 0-32 32v160h768V160z m0 288H128v416a32 32 0 0 0 32 32h704a32 32 0 0 0 32-32V448z m-256 64h128v320h-128V512z m-192 192h128v128h-128v-128z m0-192h128v128h-128v-128z m-192 192h128v128H256v-128z m0-192h128v128H256v-128z" p-id="5743" fill="#8a8a8a"></path></svg>
            </i>
            计算延迟（图片 {{ imagePageInfo.currentPage }}）</el-button>
          <!-- <el-button>打开当前截图</el-button> -->
        </el-row>
        

        <el-row justify="center" class="result-row">
          <el-col :span="4" class="info-line">
            <span>操作延迟(毫秒)</span>
          </el-col>
          <el-col :span="4" class="info-line">
            {{ result.latency}}
          </el-col>
          <el-col :span="4" class="info-line">
            <el-button link>开始图片</el-button>
          </el-col>
          <el-col :span="4" class="info-line">
            <el-button link>结束图片</el-button>
          </el-col>
        </el-row>
      </el-scrollbar>

        <el-dialog
          v-model="loadExtVideoVisible"
          title="加载外部视频"
          width="30%"
        >
          <el-input v-model="extVideoPath"></el-input>
          <template #footer>
            <span class="dialog-footer">
              <el-button @click="loadExtVideoVisible = false">取消</el-button>
              <el-button type="primary" @click="handleStartWithExtVideo"
                >确定</el-button
              >
            </span>
          </template>
        </el-dialog>
    </el-main>
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
  /* height: 80vh; */
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
.menu-tab {
  height: inherit;
}

.platform-tabs {
  height: inherit;
}

.panel-container {
  height: inherit;
}

.touch-area {
  color: rgb(46,211,111);
}

.watch-area {
  color: rgb(246,77,62);
}

.result-row {
  background-color: #ebeef0;
  margin-bottom: 7px;
}
</style>
