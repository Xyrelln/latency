<script setup lang="ts">
import {reactive, ref, inject, Ref, watch, onMounted, onUnmounted, computed} from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { ElNotification } from 'element-plus'
import {adb, core} from '../../wailsjs/go/models'
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
  GetDisplay
} from '../../wailsjs/go/app/Api'

const formSize = ref('default')
const formRef = ref<FormInstance>()
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
  sx: 0,
  sy: 0,
  dx: 0,
  dy: 0,
  speed: 500,
  interval: 2000,
  scene_id: 0,
  location: true
})

const rules = reactive<FormRules>({
  devices: [
    { required: true, message: '请选择设备', trigger: 'blur' },
  ],
  sx: [
    {
      required: true,
      message: '请输入X轴数值',
      trigger: 'blur',
    },
  ],
  sy: [
    {
      required: true,
      message: '请输入Y轴数值',
      trigger: 'blur',
    },
  ],
  dx: [
    {
      required: true,
      message: '请输入X轴数值',
      trigger: 'blur',
    },
  ],
  dy: [
    {
      required: true,
      message: '请输入Y轴数值',
      trigger: 'blur',
    },
  ],
  speed: [
    {
      required: true,
      message: '请输入移动速率（毫秒）',
      trigger: 'blur',
    },
  ],
  interval: [
    {
      required: true,
      message: '请输入间隔时间（毫秒）',
      trigger: 'blur',
    },
  ],
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      console.log('submit!')
    } else {
      console.log('error submit!', fields)
    }
  })
}

const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.resetFields()
}

const options = Array.from({ length: 10000 }).map((_, idx) => ({
  value: `${idx + 1}`,
  label: `${idx + 1}`,
}))
    

onMounted(()=>{
})

onUnmounted(()=>{
})


defineExpose({
})

// const data: Array<adb.Device> = ref([])

const data: {devices: Array<adb.Device>} = reactive({
  devices: [],
})

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


const status = ref(0)

function getDeviceList (value: any) {
  ListDevices().then(result => {
    if (result != null) {
      data.devices = result
    }
  })
}

async function handleGetDisplay() {
  await GetDisplay(form.devices[0]).then((res: adb.Display) => {
    console.log(res)
    deviceInfo.width = res.width
    deviceInfo.height = res.height
  })
}

async function setPointerLocationOn():Promise<Boolean> {
  let result = false
  await SetPointerLocationOn(form.devices[0]).then(res =>{ 
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
  SetPointerLocationOff(form.devices[0]).then(res =>{ 
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


async function handleStartRun() {
  // await handleGetDisplay()
  await setPointerLocationOn()
  const swipeEvent = adb.SwipeEvent.createFrom(
    { 
      sx: form.sx,
      sy: form.sy,
      dx: form.dx,
      dy: form.dy,
      speed: form.speed
    }
  )
  // const interval = 2
  console.log(swipeEvent)
  SetAutoSwipeOn(swipeEvent, form.interval)
  status.value = 1
}


async function handleStopRun() {
  SetAutoSwipeOff()
  status.value = 0
  setPointerLocationOff()
}


const scenes = ref([
  { id: 1, name: '第五人格-视角滑动'}
])

async function handleSceneChange(value: string) {
  console.log(value)
  await handleGetDisplay()
  form.sx = deviceInfo.height/2
  form.sy =  deviceInfo.width/2
  form.dx = deviceInfo.height/2 + deviceInfo.height/2/2
  form.dy = deviceInfo.width/2
  form.speed = 500
  form.interval = 2000
}

function handleDeviceChange(value: string) {
  console.log(value)
}
    
    
</script>

<template>
  <div>
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      class="demo-form"
      :size="formSize"
      status-icon
      >
      <el-form-item label="设备" prop="devices">
        <el-select
          v-model="form.devices"
          multiple
          @visible-change="getDeviceList"
          filterable
          @change="handleDeviceChange"
          placeholder="请选择设备"
          style="width:100%"
        >
          <el-option
            v-for="item in data.devices"
            :key="item.Serial"
            :label="item.Serial"
            :value="item.Serial"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="场景" prop="scene">
        <el-select
            v-model="form.scene_id"
            filterable
            placeholder="请选择场景"
            @change="handleSceneChange"
            style="width:100%">
          <el-option
            v-for="item in scenes"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="开始坐标" required>
        <el-col :span="8">
          <el-form-item prop="sx">
            <el-input v-model.number="form.sx" placeholder="X轴"></el-input>
          </el-form-item>
        </el-col>
        <el-col class="text-center" :span="1">
          <span class="text-gray-500">-</span>
          </el-col>
        <el-col :span="8">
          <el-form-item prop="sy">
            <el-input v-model.number="form.sy" placeholder="Y轴"></el-input>
          </el-form-item>
        </el-col>
      </el-form-item>
      <el-form-item label="结束坐标" required>
        <el-col :span="8">
          <el-form-item prop="sx">
            <el-input v-model.number="form.dx" placeholder="X轴"></el-input>
          </el-form-item>
        </el-col>
        <el-col class="text-center" :span="1">
          <span class="text-gray-500">-</span>
        </el-col>
        <el-col :span="8">
          <el-form-item prop="sx">
            <el-input v-model.number="form.dy" placeholder="Y轴"></el-input>
          </el-form-item>
        </el-col>
      </el-form-item>
      <el-form-item label="移动速率" prop="speed">
        <el-col :span="8">
          <el-input v-model.number="form.speed" placeholder="速率（毫秒）"></el-input>
        </el-col>
      </el-form-item>
      <el-form-item label="间隔时间" prop="interval">
        <el-col :span="8">
          <el-input v-model.number="form.interval" placeholder="间隔时间(毫秒)"></el-input>
        </el-col>
      </el-form-item>
      <el-form-item label="开启指针位置" prop="location">
        <el-col :span="8">
          <!-- <el-input v-model.number="form.interval" placeholder="间隔时间(毫秒)"></el-input> -->
          <!-- <el-checkbox label="开启指针位置" name="location" /> -->
          <el-switch v-model="form.location" />
        </el-col>
      </el-form-item>
      <el-form-item>
          <el-button type="primary" v-if="status===0" @click="handleStartRun">开始</el-button>
          <el-button type="danger" v-if="status===1" @click="handleStopRun">停止</el-button>
      </el-form-item>
      </el-form>

      <!-- <div v-for="d of data.devices">

      </div> -->

      <!-- <el-form>
        <el-form-item label="开始坐标" prop="region">
        <el-col :span="8">
          <el-input v-model.number="form.sx" placeholder="X轴"></el-input>
        </el-col>
        <el-col class="text-center" :span="1">
          <span class="text-gray-500">-</span>
          </el-col>
        <el-col :span="8">
          <el-input v-model.number="form.sy" placeholder="Y轴"></el-input>
        </el-col>
      </el-form-item>
      <el-form-item label="结束坐标" prop="count">
        <el-col :span="8">
          <el-input v-model.number="form.dx" placeholder="X轴"></el-input>
        </el-col>
        <el-col class="text-center" :span="1">
          <span class="text-gray-500">-</span>
        </el-col>
        <el-col :span="8">
          <el-input v-model.number="form.dy" placeholder="Y轴"></el-input>
        </el-col>
      </el-form-item>
      <el-form-item label="移动速率" required>
        <el-col :span="8">
          <el-input v-model.number="form.speed" placeholder="速率（毫秒）"></el-input>
        </el-col>
      </el-form-item>
      <el-form-item label="间隔时间" prop="type">
        <el-col :span="8">
          <el-input v-model.number="form.interval" placeholder="间隔时间(毫秒)"></el-input>
        </el-col>
      </el-form-item>
      <el-form-item>
          <el-button type="primary" @click="submitForm(formRef)">开始</el-button>
      </el-form-item> -->
    <!-- </el-form> -->

    </div>
</template>

<style>


</style>