<script setup lang="ts">
import {reactive, ref, inject, Ref, watch, onMounted, onUnmounted, computed} from 'vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
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
  ListRecords,
  UploadFile,
  // IsAppReady2,
} from '../../wailsjs/go/app/Api'
import {adb, core, fs} from '../../wailsjs/go/models'

// const files = reactive({value: []})
// const count = ref(0)
// const load = () => {
//   count.value += 2
// }

const data = ref<Array<fs.RecordFile>>([])

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

const handleLoadCacheFiles = async() => {
  ListRecords().then(res => {
    data.value = res
  }).catch(err => {
    data.value = []
    ElNotification({
      title: '数据加载失败',
      type: 'error',
      message: '缓存数据预加载失败, error: ' + err,
    })
  })
}

const hanleUploadFile = async(filePath: string) => {
  UploadFile(filePath).then(res => {
    ElNotification({
      title: '操作成功',
      type: 'success',
      message: '文件上传成功',
    })
  }).catch(err => {
    ElNotification({
      title: '操作失败',
      type: 'error',
      message: '文件上传失败, error: ' + err,
    })
  })
}


defineExpose({
  handleLoadCacheFiles,
})

</script>

<template>
  <div>
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
          content="上传所有数据"
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

    <el-scrollbar style="height: 50vh;">
      <div class="infinite-list">
        <!-- <div v-for="i in count" :key="i" class="infinite-list-item">
          202209271606{{ i }}  200 待上传
        </div> -->
        <div v-for="item in data" :key="item.dir_name" class="infinite-list-item">
          {{ item.dir_name }}({{ Math.floor(item.size/1000/1000) }}MB) <el-icon @click="hanleUploadFile(item.file_path)"><UploadFilled /></el-icon>
        </div>
      </div>

    </el-scrollbar>
  </div>
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
  opacity: 0.6;
  /* background: var(--el-color-primary-light-9); */
  margin: 7px 0;
  /* color: var(--el-color-primary); */
}
.infinite-list .infinite-list-item + .list-item {
  /* margin-top: 10px; */
}

.infinite-list-item:hover {
  background: var(--el-color-primary-light-9);
  opacity: 0.8;
}

.el-tabs__header {
  margin: 0 0 7px !important;
}

</style>
    