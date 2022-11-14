<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { UserFilled, VideoPlay, Delete } from '@element-plus/icons-vue'
import { ElMessage, TabsPaneContext } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import ScreenPreview from '../components/ScreenPreview.vue';

import { isWailsRun } from '@/utils/utils'

import { 
  ListDevices,
  InputSwipe,
  InputTap,
  LoadScreenshot,
  ListScens,
  SetScene,
  DeleteScene,
  SetPointerLocationOn,
} from '../../wailsjs/go/app/Api'
import {adb, app, core} from '../../wailsjs/go/models'

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

const scalePercent = ref(100)

const latencyForm = reactive({
  // serial: '',
  auto: true,
  scene: '',
  scenes: [],
  device: {
    serial: '',
    device: '',
  },
})

const tabName = ref('record')
const imagePreviewRef = ref()
const inputSceneName = ref('')
const deviceInfo = reactive({
  width: 1080,
  height: 1920,
})

const previewImageInfo = reactive({
  width: 0,
  height: 0,
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


const cropTouchInfo:CropArea = reactive({
  top: 0,
  left: 0,
  width: 100,
  height: 35,
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
    top: Math.trunc(cropInfo.top * imageZoom.value),
    left: Math.trunc(cropInfo.left * imageZoom.value),
    width: Math.trunc(cropInfo.width * imageZoom.value),
    height: Math.trunc(cropInfo.height * imageZoom.value),
  }
})

const realCropTouchInfo:Ref<CropArea> = computed(()=> {
  return {
    top: Math.trunc(cropTouchInfo.top * imageZoom.value),
    left: Math.trunc(cropTouchInfo.left * imageZoom.value),
    width: Math.trunc(cropTouchInfo.width * imageZoom.value),
    height: Math.trunc(cropTouchInfo.height * imageZoom.value),
  }
})

const realUserAction = computed(() => {
  return {
    // type: 'click',
    x: Math.trunc(userAction.x * imageZoom.value),
    y: Math.trunc(userAction.y * imageZoom.value),
    tx: Math.trunc(userAction.tx * imageZoom.value),
    ty: Math.trunc(userAction.ty * imageZoom.value),
    // speed: userAction.speed,
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

const userScenes: {scens: Array<app.UserScene>} = reactive({
  scens: [],
})

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
    if (d.serial == latencyForm.device.serial) {
      return d.state
    }
  }
}

/**
 * 发送操作事件
 */
 function handleInput() {
  if (userAction.type === 'swipe') {
    const swipeEvent = adb.SwipeEvent.createFrom({ 
      sx: realUserAction.value.x,
      sy: realUserAction.value.y,
      dx: realUserAction.value.tx,
      dy: realUserAction.value.ty,
      speed: userAction.speed
    })  
    InputSwipe(latencyForm.device.serial, swipeEvent)
  } else {
    const tapEvent = adb.TapEvent.createFrom({
      x: realUserAction.value.x,
      y: realUserAction.value.y
    })
    InputTap(latencyForm.device.serial, tapEvent).then().catch(err => { console.log(err)})
  }
}


const handleCropChange = (res: CropInfo)=> {
  cropInfo.left = res.left
  cropInfo.top = res.top
  cropInfo.width = res.width
  cropInfo.height = res.height
}

const handleCropTouchChange = (res: CropInfo)=> {
  cropTouchInfo.left = res.left
  cropTouchInfo.top = res.top
  cropTouchInfo.width = res.width
  cropTouchInfo.height = res.height
}

const handleGetPreviewImgSize = () => {
  const pImgSize = imagePreviewRef.value.getPreviewImgSize()
  previewImageInfo.width = pImgSize.width
  previewImageInfo.height = pImgSize.height

  imageZoom.value = imageInfo.width / pImgSize.width
}

const handleScaleChange = () => {
  // const pImgSize = imagePreviewRef.value.getPreviewImgSize()

  // console.log(pImgSize)
  // imageZoom.value = imageInfo.width / pImgSize.width
  // console.log(imageZoom.value)

  handleGetPreviewImgSize()

  imagePreviewRef.value.updateSelectBoxStyle()
  imagePreviewRef.value.updateSelectBoxTouchStyle()
}

const handlePageChange = (val: number) => {
  imagePageInfo.currentPage = val
  handleLoadImage(imagePageInfo.currentPage -1)
}

const handleOpenFolder = (val: number) => {
  // OpenImageInExplorer(val).then().catch(err => console.log(err))
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

const handleUserAction = (action: any) => {
  // console.log(action)
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

const isVerticalScreen = () => {
  return imageInfo.width < imageInfo.height
}

const handleLoadScreenshot = async () => {
  NProgress.start()
  await setPointerLocationOn()
  LoadScreenshot(latencyForm.device.serial).then((res:core.ImageInfo) => {
    imageInfo.path = res.path
    imageInfo.width = res.width
    imageInfo.height = res.height

    if (isVerticalScreen()) {
      imagePreviewRef.value.setScalePercent(50)
    } else {
      imagePreviewRef.value.setScalePercent(100)
    }

    // const pImgSize = imagePreviewRef.value.getPreviewImgSize()
    handleGetPreviewImgSize()
    // imageZoom.value = imageInfo.width / pImgSize.width

    imagePreviewRef.value.updateSelectBoxStyle()
    imagePreviewRef.value.updateSelectBoxTouchStyle()
    imagePreviewRef.value.switchSelectBoxShow(true)
    imagePreviewRef.value.switchSelectBoxTouchShow(true)

    NProgress.done()
  }).catch(err => {
    console.log(err)
    NProgress.done()
  })
}


const handleGetScenes = () => {
  ListScens().then((res: Array<app.UserScene>) => {
    console.log(res)
    userScenes.scens = res
  }).catch(err => {
    console.log(err)
  })
}

const handleSceneChange = (val:app.UserScene) => {
}

const handleSetScene = () => {
  if (inputSceneName.value === '') {
    ElMessage({
      type: 'warning',
      message: '请输入场景名称'
    })
    return
  }

  const deviceInfo = app.DeviceInfo.createFrom({
    device_name:  latencyForm.device.device,
    screen_width: imageInfo.width,
    screen_height: imageInfo.height,
  })
  const cropInfo = app.CropInfo.createFrom({
    top: realCropInfo.value.top,
    left: realCropInfo.value.left,
    width: realCropInfo.value.width,
    height: realCropInfo.value.height,
  })
  const cropTouchInfo = app.CropInfo.createFrom({
    top: realCropTouchInfo.value.top,
    left: realCropTouchInfo.value.left,
    width: realCropTouchInfo.value.width,
    height: realCropTouchInfo.value.height,
  })
  const action = app.UserAction.createFrom({
    x: realUserAction.value.x,
    y: realUserAction.value.y,
    tx: realUserAction.value.tx,
    ty: realUserAction.value.ty,
    speed: userAction.speed,
    type: userAction.type,
  })
  const userScene = app.UserScene.createFrom({
    name: inputSceneName.value,
    device: deviceInfo,
    crop_coordinate: cropInfo,
    crop_touch_coordinate: cropTouchInfo,
    action: action
  })
  SetScene(userScene).then(res => {
    inputSceneName.value = ''
    ElMessage({
      type: 'success',
      message: '场景保存成功'
    })
  }).catch(err => { console.log(err) })
}

const handleDelScene = (key: string) => {
  DeleteScene(key).then( res => {
    handleGetScenes()
  }).catch(err => { console.log(err)})
}

const handleTabClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event)
  if (tab.props.name === 'manage') {
    handleGetScenes()
  }
}


// onMounted(()=> {
//   // 如果是在 wails 运行环境则运行环境检查及事件监听
//   if (isWailsRun()) {
//   }
// })

// onUnmounted(()=>{
//   if (isWailsRun()) {
//   }
// })


</script>

<template>
  <div class="automatin-main">
   
        <el-tabs 
          v-model="tabName"
          class="platform-tabs"
          @tab-click="handleTabClick"
          >
          <el-tab-pane label="录制" name="record">
            <el-row justify="space-between" style="display:flex">
              <el-col class="button-row" :span="12">
                <el-select
                    v-model="latencyForm.device"
                    class="device-select"
                    @focus="getDeviceList"
                    filterable
                    placeholder="请选择设备"
                    >
                    <el-option
                      v-for="item in data.devices"
                      :key="item.serial"
                      :label="item.serial"
                      :value="item"
                    >
                  </el-option>
                </el-select>
                <el-button type="primary" @click="handleLoadScreenshot">
                  <el-icon>
                    <svg t="1666852111201" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6878" width="200" height="200"><path d="M817.926086 824.102251l-33.619702-33.623795c67.307965-68.77948 108.87057-162.728198 108.87057-266.540565 0-210.568786-170.671101-381.177466-381.177466-381.177466-26.38595 0-52.001351 3.009544-76.785827 8.262174l51.617611 43.2265c10.053983 8.518001 11.334139 23.503297 2.817162 33.621748-8.390088 10.05603-23.438828 11.336186-33.493835 2.883677l-106.374726-89.211834-3.136433-2.625804c-0.89744-0.704035-1.281179-1.727341-1.986237-2.560312l-4.930289-4.931313c0.255827-0.128937 0.576121-0.25685 0.896416-0.38374-2.113127-4.227278-2.818185-9.028631-2.177596-13.641695-0.063445-6.852058 2.49789-13.639648 8.198729-18.443048l109.511159-91.901083c10.055007-8.452509 25.103748-7.172353 33.493835 2.946099 8.516977 10.118452 7.235798 25.169239-2.882654 33.621748l-76.914764 64.48978c32.789801-8.196683 66.796312-13.000082 102.147448-13.000082 236.8258 0 428.824649 191.998849 428.824649 428.824649C940.824137 640.877229 893.879965 746.802724 817.926086 824.102251L817.926086 824.102251zM130.823046 523.93789c0 210.506365 170.670078 381.175419 381.176442 381.175419 40.795125 0 80.053244-6.53074 116.875894-18.442024l-67.883063-56.933687c-10.05603-8.390088-11.399631-23.438828-2.947122-33.493835 8.516977-10.053983 23.50432-11.398608 33.55728-2.947122l109.513206 91.902106c6.083555 5.122671 8.580422 12.617365 8.005324 19.981077 0.575098 7.363711-1.921769 14.793938-8.005324 19.97903l-109.447714 91.837638c-10.05603 8.388041-25.105794 7.106862-33.559327-3.009544-8.452509-9.992585-7.108908-25.106818 2.947122-33.494858l46.174646-38.679951c-30.676673 7.043417-62.441121 11.014868-95.166453 11.014868-236.891292-0.063445-428.888094-192.062294-428.888094-428.888094 0-123.023918 52.064796-233.624898 135.064139-311.819818l33.620725 33.685193C177.508321 315.354318 130.823046 414.106436 130.823046 523.93789L130.823046 523.93789z" p-id="6879" fill="#707070"></path></svg>
                  </el-icon>
                  同步屏幕
                </el-button>

                <el-button @click="handleInput" :icon="VideoPlay">
                  演示
                </el-button>
              </el-col>

              <el-col class="button-row" :span="12" >
                <div style="float:right">
                  <el-input class="scene-name-input" v-model="inputSceneName" clearable placeholder="请输入保存的场景名称"></el-input>
                  <el-button @click="handleSetScene" type="success">
                    <i class="el-icon button-icon">
                      <svg t="1666605626827" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4819" width="200" height="200"><path d="M665.6 332.8 665.6 128C665.6 113.86151 654.13849 102.4 640 102.4 625.86151 102.4 614.4 113.86151 614.4 128L614.4 332.8C614.4 346.93849 625.86151 358.4 640 358.4 654.13849 358.4 665.6 346.93849 665.6 332.8L665.6 332.8ZM640 51.2 819.2 51.2 793.6 25.6 793.6 384.133545C793.6 426.406699 759.102946 460.8 716.727898 460.8L281.672102 460.8C239.236715 460.8 204.8 426.413438 204.8 384.133545L204.8 25.6C204.8 11.46151 193.33849 0 179.2 0 165.06151 0 153.6 11.46151 153.6 25.6L153.6 384.133545C153.6 454.707134 210.976425 512 281.672102 512L716.727898 512C787.345461 512 844.8 454.718257 844.8 384.133545L844.8 25.6 844.8 0 819.2 0 640 0C625.86151 0 614.4 11.46151 614.4 25.6 614.4 39.73849 625.86151 51.2 640 51.2L640 51.2Z" p-id="4820" fill="#8a8a8a"></path><path d="M844.8 972.8 128.081132 972.8C85.544157 972.8 51.2 938.575806 51.2 896.163853L51.2 100.711064 51.2 25.6 25.6 51.2 102.4 51.2 896.233363 51.2C938.580175 51.2 972.8 85.414085 972.8 127.868001L972.8 998.4C972.8 1012.53849 984.26151 1024 998.4 1024 1012.53849 1024 1024 1012.53849 1024 998.4L1024 127.868001C1024 57.135182 966.85523 0 896.233363 0L102.4 0 25.6 0 0 0 0 25.6 0 100.711064 0 896.163853C0 966.892966 57.307204 1024 128.081132 1024L844.8 1024C858.93849 1024 870.4 1012.53849 870.4 998.4 870.4 984.26151 858.93849 972.8 844.8 972.8L844.8 972.8Z" p-id="4821" fill="#8a8a8a"></path></svg>
                    </i>
                    保存
                  </el-button>
                </div>
              </el-col>

            </el-row>

            <!-- <el-container>
              <el-main class="main-content"> -->
                  <el-row class="result-row">
                      <el-col :span="8" class="info-line">
                        <span>动作: </span>
                        <span v-if="userAction.type==='click'">点击 x: {{ realUserAction.x}} y: {{ realUserAction.y }}</span>
                        <span v-if="userAction.type==='swipe'">滑动 x: {{ realUserAction.x}} y:{{ realUserAction.y}} tx: {{ realUserAction.tx }} ty:{{ realUserAction.ty }}  s: {{ userAction.speed }}</span>
                      </el-col>
                      <el-col :span="8" class="info-line">
                        <span class="touch-area">触控区域:</span>
                        x: {{ realCropTouchInfo.left }} y: {{ realCropTouchInfo.top }}  w: {{ realCropTouchInfo.width }}  h: {{ realCropTouchInfo.height }}
                      </el-col>
                      <el-col :span="8" class="info-line">
                        <span class="watch-area">观察区域</span>
                        x: {{ realCropInfo.left }} y: {{ realCropInfo.top }}  w: {{ realCropInfo.width }}  h: {{ realCropInfo.height }}
                      </el-col>
                  </el-row>
                  <el-scrollbar style="height: calc(100vh - 216px)">
                    <!-- <div> -->
                      <ScreenPreview
                        ref="imagePreviewRef"
                        :imageInfo="imageInfo"
                        :cropInfo="cropInfo"
                        :cropTouchInfo="cropTouchInfo"
                        :pageInfo="imagePageInfo"
                        @crop-change="handleCropChange"
                        @page-change="handlePageChange"
                        @open-folder="handleOpenFolder"
                        @user-action="handleUserAction"
                        @crop-touch-change="handleCropTouchChange"
                        @scale-change="handleScaleChange"
                        />
                    <!-- </div> -->
                </el-scrollbar>
              <!-- </el-main>
            </el-container> -->
          </el-tab-pane>

          <el-tab-pane label="管理" name="manage">
            <el-scrollbar style="height: calc(100vh - 100px);width: calc(100vw - 100px)">
              <el-row v-for="(item, index) in userScenes.scens" key="item">
                {{ item }}
                <el-button type="danger" :icon="Delete" @click="handleDelScene(item.key)">删除</el-button>
              </el-row>
            </el-scrollbar>
          </el-tab-pane>
        </el-tabs>
  </div>
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
  /* margin-left: 1rem; */
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
  font-size: 11px;
  line-height: 18px;
  display:table-cell;
  vertical-align:middle;
  padding: 3px;
}

.button-row {
  margin: 7px 0;
}

.device-select {
  width: 180px;
  margin-right: 12px;
}
.scene-name-input {
  width: 180px;
  margin-right: 12px;
}

.platform-tabs {
  width: calc(100vw - 50px);
}

.automatin-main {
  /* height: calc(100vh - 88px); */
  /* margin: 0 7px; */
  padding: 7px;
  border: 1px solid #cbd5e0;
  border-radius: 4px;
}

.describ-text {
  display: inline-block;
  width: 100%;
  height: 120px;
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
