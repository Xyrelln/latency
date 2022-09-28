<script setup lang="ts">
import {reactive, ref, inject, Ref, watch, onMounted, onUnmounted, computed} from 'vue'
import { ElMessage } from 'element-plus'
import { 
  ListDevices,
  Start,
  StartRecord,
  // StopRecord,
  // StopProcessing,
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
  // IsAppReady2,
} from '../../wailsjs/go/app/Api'

const files = reactive({value: []})
const count = ref(0)
const load = () => {
  count.value += 2
}

/**
 * 清理所有缓存数据
 */
const handleClearCacheData = ()=> {
  ClearCacheData().then(res => {
    ElMessage({
      type: 'success',
      message: '清理成功'
    })
  }).catch(err => {
    ElMessage({
      type: 'error',
      message: '清理失败: ' + err
    })
  })
}

/**
 * 加载外部视频
 */
const handleLoadVideo = () => {

}

/**
 * 上传所有录制视频
 */
const handleUploadVideo = () => {

}

</script>

<template>
  <el-row justify="center">
    <el-button-group>
      <el-tooltip
        effect="dark"
        content="加载外部视频"
        placement="bottom-start"
      >
        <el-button>加载</el-button>
      </el-tooltip>

      <el-tooltip
        effect="dark"
        content="加载外部视频"
        placement="bottom-start"
      >
        <el-button>上传</el-button>
      </el-tooltip>

      <el-tooltip
        effect="dark"
        content="清理所有数据"
        placement="bottom-start"
      >
      <el-button @click="handleClearCacheData">清理</el-button>
      </el-tooltip>
    </el-button-group>
  </el-row>

    <el-scrollbar style="height: 52vh;">
      <div v-infinite-scroll="load" class="infinite-list">
        <div v-for="i in count" :key="i" class="infinite-list-item">
          202209271606{{ i }}  200 待上传
        </div>
      </div>

    </el-scrollbar>
</template>

<style>
.infinite-list {
  /* height: 100%; */
  padding: 0;
  margin: 0;
  list-style: none;
}
.infinite-list .infinite-list-item {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 32px;
  background: var(--el-color-primary-light-9);
  margin: 10px 0;
  color: var(--el-color-primary);
}
.infinite-list .infinite-list-item + .list-item {
  margin-top: 10px;
}

</style>
    