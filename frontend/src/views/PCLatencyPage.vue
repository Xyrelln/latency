<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import ScreenPreview from '../components/ScreenPreview.vue';
import FileRecord from '../components/FileRecord.vue';
import HelpPage from '../components/HelpPage.vue';

import { isWailsRun } from '@/utils/utils'

import { 
  ClearCacheData,
  IsAppReady,
  StartWinOpLatency,
  CalculateLatencyByImageDiff,
  CalculateLatencyByCurrentImage,
  GetImage,
  OpenImageInExplorer,
  ListCaptureWindows,
} from '../../wailsjs/go/app/Api'
import {adb, app, core, latencywin} from '../../wailsjs/go/models'
import {
  EventsOn,
  EventsOff,
  WindowReload,
} from '../../wailsjs/runtime/runtime'


const latencyTabName = ref('list')
const fileRecordRef = ref()
const isRunning = ref(false)

const latencyForm = reactive({
  capture_window: '',
  operate_method: 'keyboard',
  operate_key: 'KeyA',
  auto: false,
  start_hotkey: 'F2',
  diffScore: 20,
  frame_count: 300,
  auto_upload: false,
})

const capture_windows = reactive({value: []})

const cropInfo:CropArea = reactive({
  top: 50,
  left: 50,
  width: 90,
  height: 90,
})

const imagePageInfo:ImagePage = reactive({
  size: 1,
  total: 0,
  currentPage: 1,
})

const operateMethods = [
  { name: '键盘', value: 'keyboard'},
  { name: '鼠标', value: 'mouse'},
]

const imagePreviewRef = ref()

const imageInfo = reactive({
  path: '',
  width: 0,
  height: 0,
  size: 0,
  count: 0,
  index: 0,
  createTime: 0,
})

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
  frame_count: [
    {
      required: true,
      message: '截图总数',
      trigger: 'blur',
    },
    {
      validator: checkGreaterThanZero,
      trigger: 'blur',
    }
  ],
}

const result = reactive({
  latency: 0,
  responseIndex: 0,
  responseTime: 0,
  imageCount: 0,
  inputTime: 0,
  // currentImageIndex: 0,
})

function checkGreaterThanZero (rule: any, value: any, callback: any)  {
  if (value <= 0) {
    callback(new Error('数值必须大于0'))
  } else {
    callback()
  }
}

const handleLoadImage = async(imageIndex: number) => {
  GetImage(imageIndex).then((res:app.GetImageResp) => {
      console.log(res)
      imageInfo.path = res.imageFilePath
      imageInfo.width = res.imageWidth
      imageInfo.height = res.imageHeight
      imageInfo.count = res.length
      // imageInfo.createTime = res.screenshotTime
      imageInfo.index = res.currentIndex

    }).catch(err => {
      console.log(err)
    })
}

const resetStatus = () => {
  if (NProgress.isStarted()) {
    NProgress.done()
  }
  isRunning.value = false
}


/**
 * 绑定监听
 */
async function addEventLister() {
  console.log("addEventLister")
  EventsOn("latencyWindowsComplete", (res)=>{
    console.log("latencyWindowsComplete")
    console.log(res)

    // result.imageCount = res.imageCount
    result.inputTime = res.inputTime
    imagePageInfo.total = res.imageCount

    resetStatus()

    ElNotification({
      title: '进度提醒-录制完成',
      type: 'success',
      message: "录制成功",
    })


    console.log("handleLoadImage")
    const firstImageIndex = 0
    handleLoadImage(firstImageIndex)
    // result.currentImageIndex = firstImageIndex
    imagePageInfo.currentPage = 1
  })

  EventsOn("latencyWindowsMessage", (res) => {
    ElNotification({
      title: '处理过程提醒',
      type: 'info',
      message: res.message,
    })
  })

  EventsOn("latencyWindowsError", (res) => {
    resetStatus()

    ElNotification({
      title: '处理过程异常',
      type: 'error',
      message: res.error,
    })
  })
 
}

const handleStart = () => {
  console.log("handleStart")
  NProgress.start()
  isRunning.value = true
  const input_config = latencywin.InputConf.createFrom({
    type: latencyForm.operate_method,
    isAuto: latencyForm.auto,
    keyTap: latencyForm.operate_key
  })
  
  const config = latencywin.Config.createFrom({
    inputConf: input_config,
    captureWindow: latencyForm.capture_window,
    frames: latencyForm.frame_count,
    startKey: latencyForm.start_hotkey,
  })

  console.log(config)
  StartWinOpLatency(config).then(res => {
    console.log("StartWinOpLatency")
  }).catch(err => {
    console.log(err)
  })
  
}


async function removeEventLister() {
  EventsOff("latencyWindowsComplete")
  EventsOff("latencyWindowsMessage")
  EventsOff("latencyWindowsError")
}

// /**
//  * 环境检查
//  */
// async function initCheck() {
//   IsAppReady().then(res => {
//   }).catch(err => {
//     ElNotification({
//       title: '环境检查',
//       type: 'error',
//       message: err,
//       duration: 0,
//     })
//   })
// }

const handleOperateKeyFocus = (event: FocusEvent, name: string) => {

  window.onkeyup=function(e: KeyboardEvent){
    switch (name) {
      case 'operate_key':
        latencyForm.operate_key = e.code
        break
      case 'start_hotkey':
        latencyForm.start_hotkey = e.code
        break
    }

  }
  // window.onkeydown=function(e: KeyboardEvent){
  //   console.log(e)
  //   // if (unVisualKeys.indexOf(e.code) >= 0) {

  //     switch (name) {
  //       case 'operate_key':
  //         latencyForm.operate_key = e.code
  //         break
  //       case 'start_hotkey':
  //         latencyForm.start_hotkey = e.code
  //         break
  //     }
      
  //   }
  // }
}

const handleOperateKeyBlur = (event: FocusEvent, name: string) => {
  window.onkeydown = null
}

function handleReload() {
  WindowReload();
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
  OpenImageInExplorer(val).then().catch(err => console.log(err))
}

const getCaptureWindows = () => {
  ListCaptureWindows().then((res:any) => {
    console.log(res)
    capture_windows.value = res
  }).catch(err => {
    console.log(err)
  })
}

/**
 * 计算延迟
 */
const handleCalc = () => {
  console.log("handleCalcWithCurrent")
  const pImgSize = imagePreviewRef.value.getPreviewImgSize()
  const rectinfo = core.ImageRectInfo.createFrom({
    x: cropInfo.left,
    y: cropInfo.top,
    w: cropInfo.width,
    h: cropInfo.height,
    preview_width: pImgSize.width,
    preview_height: pImgSize.height,
    source_width: imageInfo.width,
    source_height: imageInfo.height,
  })
  console.log(rectinfo)
  CalculateLatencyByImageDiff(rectinfo, latencyForm.diffScore).then((res:app.WinOpLatencyResult) => {
    // console.log(res)
    result.latency = res.latency
    result.responseIndex = res.responseIndex
    result.responseTime = res.responseTime

    imagePageInfo.currentPage = result.responseIndex +1

    // 加载至目标图片
    handleLoadImage(result.responseIndex)

  }).catch(err => {
    console.log(err)
  })
}

/**
 * 根据图片计算延迟
 */
const handleCalcWithCurrent = () => {
  CalculateLatencyByCurrentImage(imageInfo.index).then((res:app.WinOpLatencyResult) => {
    console.log(res)
    result.latency = res.latency
    result.responseIndex = res.responseIndex
    result.responseTime = res.responseTime
  }).catch(err => {
    console.log(err)
  })
}

onMounted(()=> {
  // 如果是在 wails 运行环境则运行环境检查及事件监听
  if (isWailsRun()) {
    // initCheck()
    addEventLister()
  }
})

onUnmounted(()=>{
  if (isWailsRun()) {
    removeEventLister()
  }
})

</script>

<template>
    <el-scrollbar style="height: calc(100vh - 100px);width: calc(100vw - 60px)">
        <el-container>
        <el-aside class="aside-content" width="240px">
            <el-row class="row-item">
              <el-form :model="latencyForm">
                <el-form-item label="录制窗口">
                  <el-select 
                    v-model="latencyForm.capture_window" class="m-2" 
                    @focus="getCaptureWindows"
                    filterable
                    placeholder="请选择录制窗口"
                    size="default">
                    <el-option
                      v-for="item in capture_windows.value"
                      :key="item"
                      :label="item"
                      :value="item"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item label="操控方式">
                  <el-select v-model="latencyForm.operate_method" class="m-2" placeholder="请选择操控方式" size="default">
                    <el-option
                      v-for="item in operateMethods"
                      :key="item.value"
                      :label="item.name"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item v-if="latencyForm.operate_method==='keyboard'" label="操控按键">
                  <el-input 
                    v-model="latencyForm.operate_key"
                    @focus="event => handleOperateKeyFocus(event, 'operate_key')"
                    @blur="event => handleOperateKeyBlur(event, 'operate_key')"
                    placeholder="请进行按键操作">
                  </el-input>
                </el-form-item>
                <el-form-item v-if="latencyForm.operate_method==='keyboard'" label="自动操作">
                  <el-switch v-model="latencyForm.auto" />
                </el-form-item>
                <el-form-item label="快捷启动">
                  <el-input 
                    v-model="latencyForm.start_hotkey"
                    @focus="event => handleOperateKeyFocus(event, 'start_hotkey')"
                    @blur="event => handleOperateKeyBlur(event, 'start_hotkey')"
                    placeholder="请进行按键操作"
                    >
                  </el-input>
                </el-form-item>
              </el-form>
            </el-row>
            <el-row class="row-item">
            <el-button class="operation-button" type="primary" @click="handleStart" :disabled="isRunning" >开始</el-button>
            </el-row>

            <el-tabs 
                v-model="latencyTabName" 
                class="platform-tabs">
                <el-tab-pane label="记录" name="list">
                <FileRecord ref="fileRecordRef"/>

                </el-tab-pane>
            
                <el-tab-pane label="设置" name="setting">
                <!-- <el-scrollbar style="height:60vh"> -->
                    <el-row>
                    <el-form :model="latencyForm" ref="settingFormRef" :rules="rules" label-position="left" label-width="100px">
                        <el-form-item label="对比阈值" prop="diffScore">
                          <el-input v-model.number="latencyForm.diffScore"/>
                        </el-form-item>
                        <el-form-item label="截图总数" prop="frame_count">
                          <el-input v-model.number="latencyForm.frame_count"/>
                        </el-form-item>
                        <el-form-item label="自动上传">
                          <el-switch v-model="latencyForm.auto_upload" />
                        </el-form-item>
                        <el-form-item label="调式">
                          <el-button @click="handleReload">重载页面</el-button>
                        </el-form-item>
                    </el-form>
                    </el-row>
                <!-- </el-scrollbar> -->
                </el-tab-pane>
                <!-- <el-tab-pane label="帮助" name="detail" disabled>
                    <HelpPage></HelpPage>
                </el-tab-pane> -->
                
            </el-tabs>
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
                />
            </div>
            <el-row justify="center" class="button-row">
              <el-button type="success" @click="handleCalc">
                <i class="el-icon button-icon">
                  <svg t="1666320784905" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5742" width="200" height="200"><path d="M928 1024H96a96 96 0 0 1-96-96V96a96 96 0 0 1 96-96h832a96 96 0 0 1 96 96v832a96 96 0 0 1-96 96zM896 160a32 32 0 0 0-32-32H160a32 32 0 0 0-32 32v160h768V160z m0 288H128v416a32 32 0 0 0 32 32h704a32 32 0 0 0 32-32V448z m-256 64h128v320h-128V512z m-192 192h128v128h-128v-128z m0-192h128v128h-128v-128z m-192 192h128v128H256v-128z m0-192h128v128H256v-128z" p-id="5743" fill="#8a8a8a"></path></svg>
                </i>
                计算延迟
              </el-button>
              <el-button @click="handleCalcWithCurrent">
                <i class="el-icon button-icon">
                  <svg t="1666320784905" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5742" width="200" height="200"><path d="M928 1024H96a96 96 0 0 1-96-96V96a96 96 0 0 1 96-96h832a96 96 0 0 1 96 96v832a96 96 0 0 1-96 96zM896 160a32 32 0 0 0-32-32H160a32 32 0 0 0-32 32v160h768V160z m0 288H128v416a32 32 0 0 0 32 32h704a32 32 0 0 0 32-32V448z m-256 64h128v320h-128V512z m-192 192h128v128h-128v-128z m0-192h128v128h-128v-128z m-192 192h128v128H256v-128z m0-192h128v128H256v-128z" p-id="5743" fill="#8a8a8a"></path></svg>
                </i>
                计算延迟（图片 {{ imagePageInfo.currentPage }}）</el-button>
            </el-row>
            <el-row justify="center" class="result-row">
              <el-col :span="4" class="info-line">
                <span>操作时间</span>
              </el-col>
              <el-col :span="4" class="info-line">
                {{ result.inputTime}}
              </el-col>

              <el-col :span="4" class="info-line">
                <span>响应时间</span>
              </el-col>
              <el-col :span="4" class="info-line">
                {{ result.responseTime}}
              </el-col>
            </el-row>
            <el-row justify="center" class="result-row">
              <el-col :span="4" class="info-line">
                <span>操作延迟(毫秒)</span>
              </el-col>
              <el-col :span="4" class="info-line">
                {{ result.latency}}
              </el-col>
              <el-col :span="4" class="info-line">
              </el-col>
              <el-col :span="4" class="info-line">
              </el-col>
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
