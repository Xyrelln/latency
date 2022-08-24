<script setup lang="ts">
import {reactive, ref, inject, Ref, onMounted, computed} from 'vue'
import { UserFilled } from '@element-plus/icons-vue'

import {ListDevices, StartRecord, StopRecord} from '../../wailsjs/go/main/App'
import {adb} from '../../wailsjs/go/models'

const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})


const deviceSelected = ref("")

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
      <el-button @click="handleStartRecord">录制</el-button>
      <el-button @click="handleStopRecord">停止录制</el-button>
      <el-button>转图片</el-button>
      <el-button>解析</el-button>
    </el-row>

  </div>

    
</template>

<style scoped>

</style>