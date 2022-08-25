<script setup lang="ts">
import {reactive, ref, inject, Ref, onMounted, computed} from 'vue'
import { UserFilled } from '@element-plus/icons-vue'

import {ListDevices, StartRecord, StopRecord, StopProcessing} from '../../wailsjs/go/main/App'
import {adb} from '../../wailsjs/go/models'


const deviceSelected = ref("")
const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})

const processStatus = ref(0)

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
  StartRecord(deviceSelected.value)
}

function handleStopRecord() {
  StopRecord(deviceSelected.value)
}

function handleStopProcessing() {
  StopProcessing()
}

</script>

<template>
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
      <el-button v-if="processStatus===0" type="primary" @click="processStatus=1" style="width: 100%">开始</el-button>
      <el-button v-if="processStatus===1" type="success" @click="processStatus=2" style="width: 100%">准备: 3</el-button>
      <el-button v-if="processStatus===2" type="danger"  @click="processStatus=0" style="width: 100%">停止: 10</el-button>
    </el-row>
    <el-row>
      <el-button @click="handleStartRecord">录制</el-button>
      <el-button @click="handleStopProcessing">停止处理</el-button>
      <el-button>转图片</el-button>
      <el-button>准备</el-button>
      <el-button>解析</el-button>
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
              <el-row class="info-list">
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
              </el-row>
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