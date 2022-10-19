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
  GetDisplay,
  IsAppReady,
  StartWinOpLatency,
  CalculateLatencyByImageDiff,
  CalculateLatencyByCurrentImage,
  GetImage,
} from '../../wailsjs/go/app/Api'
import {adb, core, latencywin} from '../../wailsjs/go/models'
import {
  EventsOn,
  EventsOff,
  WindowReload,
} from '../../wailsjs/runtime/runtime'

const deviceSelected = ref("")

const latencyTabName = ref('list')
const placeholder = "./src/assets/images/placeholder.png"
const fileRecordRef = ref()


const latencyForm = reactive({
  operate_method: 'keyboard',
  operate_key: 'A',
  auto: false,
  start_hotkey: 'F2',
  diffScore: 20,
  frame_count: 300,
  auto_upload: false,
})

const operateMethod = ref('keyboard')
const operateMethods = ['keyboard', 'mouse']


const imagePreviewRef = ref()
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


// provide('threshold', settingForm.diffScore)


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


function checkGreaterThanZero (rule: any, value: any, callback: any)  {
  if (value <= 0) {
    callback(new Error('数值必须大于0'))
  } else {
    callback()
  }
}


/**
 * 绑定监听
 */
async function addEventLister() {
 
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

const handleStart = () => {
  const input_config = latencywin.InputConf.createFrom({
    type: latencyForm.operate_method,
    isAuto: latencyForm.auto,
    keyTap: latencyForm.operate_key
  })
  
  const config = latencywin.Config.createFrom({
    inputCconf: input_config,
    imageDiff_threshold: latencyForm.diffScore,
    frames: latencyForm.frame_count,
    startKey: latencyForm.start_hotkey,
  })

  StartWinOpLatency(config).then(res => {

  }).catch(err => {

  })
  
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

const handleOperateKeyFocus = (event: FocusEvent) => {
  // event.stopPropagation()
  document.onkeydown=function(e){
    console.log(e.key)
    console.log(e.code)
    console.log(e)
    latencyForm.operate_key = e.key
  }
}

const handleOperateKeyBlur = (event: FocusEvent) => {
  document.onkeydown = null
}

onMounted(()=> {
  // 如果是在 wails 运行环境则运行环境检查及事件监听
  if (isWailsRun()) {
    initCheck()
    addEventLister()
    fileRecordRef.value.handleLoadCacheFiles()
  }
  // addKeydownListen()
})

function handleClearCache() {
  ClearCacheData()
}

function handleReload() {
  WindowReload();
}

function handleStopProcessing() {

}

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
    <el-scrollbar style="height: calc(100vh - 100px);width: calc(100vw - 60px)">
        <el-container>
        <el-aside class="aside-content" width="240px">
            <el-row class="row-item">
              <el-form :model="latencyForm">
                <el-form-item label="操控方式">
                  <el-select v-model="operateMethod" class="m-2" placeholder="Select" size="default">
                    <el-option
                      v-for="item in operateMethods"
                      :key="item"
                      :label="item"
                      :value="item"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item label="操控按键">
                  <el-input 
                    v-model="latencyForm.operate_key"
                    @focus="handleOperateKeyFocus"
                    @blur="handleOperateKeyBlur"
                    placeholder="请进行按键操作">
                  </el-input>
                </el-form-item>
                <el-form-item label="自动操作">
                  <el-switch v-model="latencyForm.auto" />
                </el-form-item>
                <el-form-item label="快捷启动">
                  <el-input v-model="latencyForm.start_hotkey" placeholder="请进行按键操作">
                  </el-input>
                </el-form-item>
              </el-form>
            </el-row>
            <el-row class="row-item">
            <el-button class="operation-button" type="primary" @click="handleStart" >开始</el-button>
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
                        <el-form-item label="截图帧数" prop="timeout">
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
                <el-tab-pane label="帮助" name="detail" disabled>
                    <HelpPage></HelpPage>
                </el-tab-pane>
                
            </el-tabs>
        </el-aside>
        <el-main class="main-content">
            <div>
              <ScreenPreview
              ref="imagePreviewRef"
              :data="imageInfo"
              />
            </div>
            <el-row>
              <el-button>计算延迟</el-button>
              <el-button>计算延迟（当前截图）</el-button>
              <el-button>打开当前截图</el-button>
            </el-row>
            <el-row>
              <span>操作时间</span>
              <span>响应时间</span>
              <span>操作延迟</span>
            </el-row>
            <!-- <el-row>
              <el-col :span="4">
                <el-input placeholder="left"></el-input>
              </el-col>
              <el-col :span="4">
                <el-input placeholder="top"></el-input>
              </el-col>
              <el-col :span="4">
                <el-input placeholder="width"></el-input>
              </el-col>
              <el-col :span="4">
                <el-input placeholder="height"></el-input>
              </el-col>
            </el-row> -->
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
</style>
