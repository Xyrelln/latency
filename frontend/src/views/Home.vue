<script setup lang="ts">
import {reactive, ref, inject, Ref, onMounted, computed, watch, onUnmounted} from 'vue'
import { UserFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
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
  SetPointerLocationOn
} from '../../wailsjs/go/app/Api'
import {adb} from '../../wailsjs/go/models'
import RightContent from '../components/RightContent.vue';
import ImagePreview from '../components/ImagePreview.vue';
import {
  EventsOn,
  EventsOff
} from '../../wailsjs/runtime/runtime'

const deviceSelected = ref("")
const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})

const rightContentRef = ref()
const startButtonText = ref("开始")
const interval = ref()
const processStatus = ref(0)
const developMode = ref(true)
const countDownSecond = ref(0)


const tabName = ref('detail')
const deviceInfo = reactive({
  android_version: null,
  cpu_arch: '',
  cpu_core_count: null,
  hardware: '',
  mem_total: 0,
  openGLES_version: '',
  product_model: ''
})

const settingForm = reactive({
  captureInterval: 1,
  jankThreshold: 83.3333
})

const rules =  {
  jankThreshold: [
    {
      required: true,
      message: '帧间隔必填',
      trigger: 'blur',
    },
    {
      validator: checkJankThreshold,
      trigger: 'blur',
    }
  ],
  captureInterval: [
    {
      required: true,
      message: '采集间隔必填',
      trigger: 'blur',
    },
    {
      validator: checkJankThreshold,
      trigger: 'blur',
    }
  ]
}

function checkJankThreshold (rule: any, value: any, callback: any)  {
  if (value <= 0) {
    callback(new Error('帧间隔必须大于0'))
  } else {
    callback()
  }
}


function getDeviceList () {
    ListDevices().then(result => {
    data.devices = result
  })
}


function handleStartRecord() {
  // clear first
  clearCurrentInterval()

  // const err = StartRecord(deviceSelected.value)
  Start(deviceSelected.value, 10)

  // add stop count down
  processStatus.value = 2
  runUntilCountDown(10)
}

function handleStopRecord() {
  // clear first
  clearCurrentInterval()

  StopRecord(deviceSelected.value)
  processStatus.value = 0
  SetPointerLocationOff(deviceSelected.value)
  handleToImage()
}

function handleStopProcessing() {
  StopProcessing()
}

function handleToImage() {
  rightContentRef.value.setPrepareDataPercent(5)
  StartTransform()
  rightContentRef.value.setPrepareDataPercent(100)
  rightContentRef.value.setAnaysePercent(5)
  handleImageAnalyse()
  rightContentRef.value.setAnaysePercent(100)
}

function handleImageAnalyse() {
  StartAnalyse().then((res)=>{
    rightContentRef.value.loadResponseTimeData(res)
  })
}

watch(countDownSecond, (value)=> {
  if (processStatus.value === 2) {
    rightContentRef.value.setRecordPercent(100 - value*10)
  }
})


function clearCurrentInterval() {
  if (interval.value != null) {
    clearInterval(interval.value)
    interval.value = null
  }
}

function runUntilCountDown(second: number, callback?: Function){
  countDownSecond.value = second
  function countDown() {
    if (countDownSecond.value  > 0) {
      countDownSecond.value  --
    } else {
      // clearCurrentInterval()
      if (callback) { callback() }
    }
  }

  clearCurrentInterval()
  interval.value = setInterval(countDown, 1000)
}

async function handlePrepare(){
  if (deviceSelected.value === "") {
    ElMessage({
      type: 'error',
      message: '请选择设备'
    })
    return
  }
  NProgress.start()
  const result = await setPointerLocationOn()
  if (result){
    processStatus.value = 1
    runUntilCountDown(3, handleStartRecord)
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


// function setEventsLister(callback: Function) {
//   EventsOn("latency:done", (data: any) => {
//     console.log("perfdata: " ,data)
//     callback(data)
//   })
// }


onMounted(()=> {
  EventsOn("latency:record_start", ()=>{
    console.log("record_start")
  })
  EventsOn("latency:record_filish", ()=>{
    setPointerLocationOff()
  })
  EventsOn("latency:transform_start", ()=>{
    rightContentRef.value.setPrepareDataPercent(5)
  })
  EventsOn("latency:transform_filish", ()=>{
    rightContentRef.value.setPrepareDataPercent(100)
  })
  EventsOn("latency:analyse_start", ()=>{
    rightContentRef.value.setAnaysePercent(5)
  })
  EventsOn("latency:analyse_filish", (res)=>{
    rightContentRef.value.setAnaysePercent(100)
    processStatus.value = 0
    rightContentRef.value.loadResponseTimeData(res)
    NProgress.done()
  })
  
})

onUnmounted(()=>{
  EventsOff("latency:record_start")
  EventsOff("latency:record_filish")
  EventsOff("latency:transform_start")
  EventsOff("latency:transform_filish")
  EventsOff("latency:analyse_start")
  EventsOff("latency:analyse_filish")
})

</script>

<template>
    <el-container>
      <el-aside width="200px">
        <div>
          <el-row>
            <el-avatar :icon="UserFilled" />
          </el-row>
          <el-row>
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
          </el-row>
          <el-row>
            <el-button v-if="processStatus===0" :disabled="deviceSelected===''" type="primary" @click="handlePrepare" style="width: 100%">准备</el-button>
            <el-button v-if="processStatus===1" type="success" @click="handleStartRecord" style="width: 100%">开始 {{ countDownSecond > 0 ? ": " + countDownSecond : ""}}</el-button>
            <el-button v-if="processStatus===2" type="danger"  @click="handleStopProcessing" style="width: 100%">停止 {{ countDownSecond > 0 ? ": " + countDownSecond : ""}}</el-button>
          </el-row>
          <el-row v-if="developMode">
            <el-button @click="handleStartRecord">rec</el-button>
            <el-button @click="handleStopProcessing">stop</el-button>
            <el-button @click="handleToImage">to_img</el-button>
            <el-button @click="handleImageAnalyse">ana</el-button>
            <el-button @click="setPointerLocationOn">set_pl_on</el-button>
            <el-button @click="setPointerLocationOff">set_pl_off</el-button>
          </el-row>

          <el-tabs 
              v-model="tabName" 
              class="platform-tabs">
              <el-tab-pane label="详情" name="detail">
                <el-scrollbar style="height:60vh">
                  <div>
                    <el-row class="info-list">
                      <el-col :span="12" class="info-line">
                        <span class="info-key">名称</span>
                      </el-col >
                      <el-col :span="12" class="info-line">
                        <span class="info-value">数值</span>
                      </el-col>
                    </el-row>
                    <!-- <el-row class="info-list">
                      <el-col :span="12" class="info-line">
                        <span class="info-key">OS Version</span>
                      </el-col >
                      <el-col :span="12" class="info-line">
                        <span class="info-value">{{ deviceInfo.android_version }}</span>
                      </el-col>
                    </el-row>
                    <el-row class="info-list">
                      <el-col :span="12" class="info-line">
                        <span class="info-key">CPU Arch</span>
                      </el-col >
                      <el-col :span="12" class="info-line">
                        <span class="info-value">{{ deviceInfo.cpu_arch }}</span>
                      </el-col>
                    </el-row>
                    <el-row class="info-list">
                      <el-col :span="12" class="info-line">
                        <span class="info-key">CPU Core</span>
                      </el-col >
                      <el-col :span="12" class="info-line">
                        <span class="info-value">{{ deviceInfo.cpu_core_count }}</span>
                      </el-col>
                    </el-row>
                    <el-row class="info-list">
                      <el-col :span="12" class="info-line">
                        <span class="info-key">MEM(G)</span>
                      </el-col >
                      <el-col :span="12" class="info-line">
                        <span class="info-value">{{ Math.floor(deviceInfo.mem_total/1000) }}</span>
                      </el-col>
                    </el-row>
                    <el-row class="info-list">
                      <el-col :span="12" class="info-line">
                        <span class="info-key">Model Name</span>
                      </el-col >
                      <el-col :span="12" class="info-line">
                        <span class="info-value">{{ deviceInfo.product_model }}</span>
                      </el-col>
                    </el-row>
                    <el-row class="info-list">
                      <el-col :span="12" class="info-line">
                        <span class="info-key">Hardware</span>
                      </el-col >
                      <el-col :span="12" class="info-line">
                        <span class="info-value">{{ deviceInfo.hardware }}</span>
                      </el-col>
                    </el-row> -->
                  </div>
                
                </el-scrollbar>
              </el-tab-pane>

              <el-tab-pane label="设置" name="setting">
                <el-scrollbar style="height:60vh">
                  <el-row>
                    <el-form :model="settingForm" ref="settingFormRef" :rules="rules">
                      <el-form-item label="采集间隔(s)" prop="captureInterval">
                        <el-input v-model.number="settingForm.captureInterval"/>
                      </el-form-item>
                      <el-form-item label="帧间隔(ms)" prop="jankThreshold">
                        <el-input v-model.number="settingForm.jankThreshold"/>
                      </el-form-item>
                    </el-form>
                  </el-row>
                </el-scrollbar>
              </el-tab-pane>
            </el-tabs>

        </div>
      </el-aside>
      <el-main>
        <RightContent ref="rightContentRef"/>
        <ImagePreview ref="ImagePreviewRef"/>
      </el-main>
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

</style>
