<script setup lang="ts">
import {reactive, ref, inject, Ref, watch, onMounted, onUnmounted, computed} from 'vue'
import {adb, core} from '@/../wailsjs/go/models'
import { ElMessage } from 'element-plus'

interface Props {
  imageInfo: core.ImageInfo
  cropInfo: CropInfo
  pageInfo: ImagePage
}

interface Emits {
  (e: 'crop-change', val: CropInfo): void
  (e: 'page-change', val: CropInfo): void
  (e: 'get-previous-page'): void
  (e: 'get-next-page'): void
  (e: 'page-change', val: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const defaultImageHolder = './assets/images/placeholder.png'

// element refs
const selectBoxRef = ref()
const previewImgRef = ref()
const resizeTopRef = ref()
const resizeRightRef = ref()
const resizeBottomRef = ref()
const resizeLeftRef = ref()
const isPaged = ref(true)


// const threshold = inject('threshold') as number
// const currentPage = ref(1)
// const pageSize = ref(1)
const small = ref(false)
const paginationDisabled = ref(false)
// const total = ref(10)
// const imgs = ref<Array<string>>([])
// const isImgLoaded = ref(true)

const selectBoxStyle = reactive({
  width: '446px', 
  height: '70px'
})
const calcButtonDisable = ref(true)

const delayTimes = ref<number>()
const location = reactive({
  x: 0,
  y: 0,
  w: 0,
  h: 0,
})

/**
 * 鼠标点击事件处理
 * @param e 
 */
const mouseDownHandler = function (e:any) {
  e.stopPropagation()

  // Get the current mouse position
  location.x = e.clientX;
  location.y = e.clientY;

  // Calculate the dimension of element
  const styles = window.getComputedStyle(selectBoxRef.value);
  location.w = parseInt(styles.width, 10);
  location.h = parseInt(styles.height, 10);

  // Attach the listeners to `document`
  document.addEventListener('mousemove', mouseMoveHandler);
  document.addEventListener('mouseup', mouseUpHandler);
};


/**
 * 鼠标移动事件处理
 * @param e 
 */
const mouseMoveHandler = function (e: any) {
  // console.log("mouseMoveHandler")
  // How far the mouse has been moved
  const dx = e.clientX - location.x;
  const dy = e.clientY - location.y;

  // Adjust the dimension of element
  selectBoxStyle.width = `${location.w + dx}px`;
  selectBoxStyle.height = `${location.h + dy}px`;
};

/**
 * 鼠标松开事件处理
 */
const mouseUpHandler = function () {
  // Remove the handlers of `mousemove` and `mouseup`
  document.removeEventListener('mousemove', mouseMoveHandler);
  document.removeEventListener('mouseup', mouseUpHandler);

  emit('crop-change', {
    left: selectBoxRef.value.offsetLeft - previewImgRef.value.offsetLeft, 
    top: selectBoxRef.value.offsetTop - previewImgRef.value.offsetTop, 
    width: selectBoxRef.value.offsetWidth, 
    height: selectBoxRef.value.offsetHeight, 
  })
};

/**
 * 初始化选择区域
 */
function selectBoxInit() {
  previewImgRef.value.addEventListener('load', ()=>{
    selectBoxRef.value.style.top = previewImgRef.value.offsetTop + props.cropInfo.top + 'px'
    selectBoxRef.value.style.left = previewImgRef.value.offsetLeft + props.cropInfo.left + 'px'
  })

  selectBoxRef.value.addEventListener('mousedown', (ev:any) => {
    const X = ev.clientX - ev.target.offsetLeft;
    const Y = ev.clientY - ev.target.offsetTop;

    const px = previewImgRef.value.offsetLeft
    const py = previewImgRef.value.offsetTop
    const pw = previewImgRef.value.offsetWidth
    const ph = previewImgRef.value.offsetHeight

    // 鼠标移动
    document.onmousemove = (ev:any) => {
      const tx = ev.clientX - X
      const ty = ev.clientY - Y

      const width = Number(selectBoxStyle.width.replace('px', ''))
      const height = Number(selectBoxStyle.height.replace('px', ''))
      
      if (tx <= px) {
        selectBoxRef.value.style.left = px + 'px'
       
      } else if (tx >= px + pw - width) {
        selectBoxRef.value.style.left = px + pw - width  + 'px'
      } else {
        selectBoxRef.value.style.left = tx + 'px'
      }

      if (ty <= py) {
        selectBoxRef.value.style.top = py + 'px'
      } else if (ty >= py + ph - height){
        selectBoxRef.value.style.top = py + ph - height + 'px'
      } else {
        selectBoxRef.value.style.top = ty + 'px'
      }
    }

    document.onmouseup = (ev:any) => {``
      document.onmousemove = null;

      emit('crop-change', {
        left: selectBoxRef.value.offsetLeft - previewImgRef.value.offsetLeft, 
        top: selectBoxRef.value.offsetTop - previewImgRef.value.offsetTop, 
        width: selectBoxRef.value.offsetWidth, 
        height: selectBoxRef.value.offsetHeight, 
      })
    }
    return false;
  })

  // add resize event listener
  resizeTopRef.value.addEventListener('mousedown', mouseDownHandler);
  resizeRightRef.value.addEventListener('mousedown', mouseDownHandler);
  resizeBottomRef.value.addEventListener('mousedown', mouseDownHandler);
  resizeLeftRef.value.addEventListener('mousedown', mouseDownHandler);
}

// function handleCalcCostTime() {
//   const rectinfo = core.ImageRectInfo.createFrom({
//     x: selectBoxRef.value.offsetLeft - previewImgRef.value.offsetLeft,
//     y: selectBoxRef.value.offsetTop - previewImgRef.value.offsetTop,
//     w: selectBoxRef.value.offsetWidth,
//     h: selectBoxRef.value.offsetHeight,
//     preview_width: previewImgRef.value.offsetWidth,
//     preview_height: previewImgRef.value.offsetHeight,
//     source_width: props.imageInfo.width,
//     source_height: props.imageInfo.height,
//   })

//   StartAnalyse(rectinfo, threshold)
//   NProgress.start()
//   delayTimes.value = 0 
//   calcButtonDisable.value = true
// }

const getPreviewImgSize = () => {
  const w = previewImgRef.value.offsetWidth
  const h = previewImgRef.value.offsetHeight

  return { width: w, height: h}
}

function handleImageLoadSuccess() {
  console.log("handleImageLoadSuccess")
}

function enableCalcButton() {
  calcButtonDisable.value = false
}

function setCalcButtonDisable(value: boolean) {
  calcButtonDisable.value = value
}

function setDefaultTime() {
  delayTimes.value = 0
}


// const handleSizeChange = (val: number) => {
//   console.log(`${val} items per page`)
// }

const handleCurrentChange = (val: number) => {
  console.log(`current page: ${val}`)
  emit('page-change', val)
  // imageInfo.path = imgs.value[val -1]
}

const getPreviousImage = ()=> {
  if (props.pageInfo.currentPage === 0 ) {
    ElMessage({
      type: 'warning',
      message: '当前为第一张图片'
    })
  } else {
    const page = props.pageInfo.currentPage - 1
    emit('page-change', page)
  }
  // emit('get-previous-page')
}

const getNextImage = ()=> {
  if (props.pageInfo.currentPage === props.pageInfo.total ) {
    ElMessage({
      type: 'warning',
      message: '当前为最后一张图片'
    })
  } else {
    const page = props.pageInfo.currentPage + 1
    emit('page-change', page)
  }
  // emit('page-change', val)
  // const page = props.pageInfo.currentPage + 1
  // emit('page-change', page)
  // emit('get-next-page')
}

onMounted(()=>{
  selectBoxInit()
})


defineExpose({
  enableCalcButton,
  setCalcButtonDisable,
  setDefaultTime,
  getPreviewImgSize,
})


</script>

<template>
  <div>
    <el-row justify="center" class="preview-content">
      <!-- <el-col> -->
        <!-- <span>标识检测区域</span> -->
        <!-- <div class="out-img-box"> -->
          <span v-if="isPaged" @click="getPreviousImage" class="page-button el-image-viewer__btn el-image-viewer__prev"><i class="el-icon"><svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"><path fill="currentColor" d="M609.408 149.376 277.76 489.6a32 32 0 0 0 0 44.672l331.648 340.352a29.12 29.12 0 0 0 41.728 0 30.592 30.592 0 0 0 0-42.752L339.264 511.936l311.872-319.872a30.592 30.592 0 0 0 0-42.688 29.12 29.12 0 0 0-41.728 0z"></path></svg></i></span>
          <img ref="previewImgRef" class="preview-img" draggable="false" :src="props.imageInfo.path == '' ? defaultImageHolder : props.imageInfo.path" alt=""/>
          <span v-if="isPaged" @click="getNextImage"  class="page-button el-image-viewer__btn el-image-viewer__next"><i class="el-icon"><svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"><path fill="currentColor" d="M340.864 149.312a30.592 30.592 0 0 0 0 42.752L652.736 512 340.864 831.872a30.592 30.592 0 0 0 0 42.752 29.12 29.12 0 0 0 41.728 0L714.24 534.336a32 32 0 0 0 0-44.672L382.592 149.376a29.12 29.12 0 0 0-41.728 0z"></path></svg></i></span>
          <div ref="selectBoxRef" :style="selectBoxStyle" class="s-move-content-header" id="select-box">
            <div ref="resizeTopRef" class="resizer resizer-t"></div>
            <div ref="resizeRightRef" class="resizer resizer-r"></div>
            <div ref="resizeBottomRef" class="resizer resizer-b"></div>
            <div ref="resizeLeftRef" class="resizer resizer-l"></div>
          </div>
          
      <!-- </el-col> -->
    </el-row>
    <!-- <el-row justify="center">
      <el-col :span="6"><span>left: </span></el-col>
      <el-col :span="6"><span>top: </span></el-col>
      <el-col :span="6"><span>width: </span></el-col>
      <el-col :span="6"><span>height: </span></el-col>
    </el-row> -->
    <el-row v-if="isPaged" justify="center">
        <el-pagination
          :currentPage="props.pageInfo.currentPage"
          :page-size="1"
          :small="small"
          :disabled="paginationDisabled"
          :background="false"
          layout="total, prev, pager, next, jumper"
          :total="props.pageInfo.total"
          @current-change="handleCurrentChange"
        />
    </el-row>
    <!-- <el-row justify="space-between" class="item-result">
      <el-col :span="4">
        <el-button type="primary" @click="handleCalcCostTime" :disabled="calcButtonDisable">计算延迟</el-button>
      </el-col>
      <el-col :span="12">
        <span class="data-item">结果（毫秒）: {{ delayTimes }}</span>
      </el-col>
    </el-row>
    <el-row> -->
     
    <!-- </el-row> -->
  </div>
</template>

<style>

.out-img-box {
  /* width: 100%; */
  width: 500px;
  width: 500px;
  height: 100%;
  /* line-height: 600px; */
  text-align: center;
}

.preview-img {
    max-width: 500px;
    max-height: 500px;
    vertical-align: middle;
    align-items: center;
}

#select-box {
  /* width: 200px;
  height: 200px; */
  background: rgba(255, 255, 0, 0.4);
  position: absolute;
  display: none;
  cursor: move;
  display: block;
}

img {
  -webkit-user-drag: none;
}

/* body {
  -webkit-user-select: none
} */

.resizable {
  position: relative;

  /* Center the content */
  align-items: center;
  display: flex;
  justify-content: center;

  /* Misc */
  border: 1px solid #cbd5e0;
  height: 8rem;
  width: 8rem;
}
.resizer {
    /* All resizers are positioned absolutely inside the element */
    position: absolute;
}

/* Placed at the right side */
.resizer-r {
    cursor: col-resize;
    height: 100%;
    right: 0;
    top: 0;
    width: 5px;
}

.resizer-l {
    cursor: col-resize;
    height: 100%;
    left: 0;
    top: 0;
    width: 5px;
}


/* Placed at the bottom side */
.resizer-b {
    bottom: 0;
    cursor: row-resize;
    height: 5px;
    left: 0;
    width: 100%;
}

.resizer-t {
    top: 0;
    cursor: row-resize;
    height: 5px;
    left: 0;
    width: 100%;
}


.data-item {
  display: block;
}


.demo-pagination-block + .demo-pagination-block {
  margin-top: 10px;
}
.demo-pagination-block .demonstration {
  margin-bottom: 16px;
}

.item-result {
  margin-top: 0.5rem;
}

.page-button:hover {
  color: rgb(90, 156, 248);
}


</style>