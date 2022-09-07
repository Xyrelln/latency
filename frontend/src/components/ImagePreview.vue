<script setup lang="ts">
import {reactive, ref, inject, Ref, watch, onMounted, onUnmounted, computed} from 'vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import {adb, core} from '../../wailsjs/go/models'
import { 
  StartAnalyse,
} from '../../wailsjs/go/app/Api'
import {
  EventsOn,
  EventsOff,
} from '../../wailsjs/runtime/runtime'

interface Props {
  data: core.ImageInfo
}

const imageInfo = reactive({
  path: '',
  width: 0,
  height: 0,
  size: 0,
})
const props = defineProps<Props>()
const selectBoxRef = ref()
const previewImgRef = ref()
// const canvasRef = ref()
// const resizeRef = ref()
const resizeTopRef = ref()
const resizeRightRef = ref()
const resizeBottomRef = ref()
const resizeLeftRef = ref()
const selectBoxStyle = reactive({
  width: '200px', 
  height: '200px'
})
const calcButtonDisable = ref(true)

const delayTimes = ref<Array<number>>([])
const location = reactive({
  x: 0,
  y: 0,
  w: 0,
  h: 0,
})

const selectArea = reactive({
    ax: 0,
    ay: 0,
    mx: 0,
    my: 0,
    bx: 0,
    by: 0,
    width: 0,
    height: 0,
    paint: false
})

const threshold = inject('threshold') as number

function handleImageAnalyse() {
  const rectinfo = core.ImageRectInfo.createFrom({

  })
  StartAnalyse(rectinfo, threshold).then((res)=>{
  })
}

function selectBoxInit() {
  previewImgRef.value.addEventListener('load', ()=>{
    selectBoxRef.value.style.top = previewImgRef.value.offsetTop + 'px'
    selectBoxRef.value.style.left = previewImgRef.value.offsetLeft + 'px'
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

    document.onmouseup = (ev:any) => {
      document.onmousemove = null;
    }
    return false;
  })

}

onMounted(()=>{
  selectBoxInit()

  resizeTopRef.value.addEventListener('mousedown', mouseDownHandler);
  resizeRightRef.value.addEventListener('mousedown', mouseDownHandler);
  resizeBottomRef.value.addEventListener('mousedown', mouseDownHandler);
  resizeLeftRef.value.addEventListener('mousedown', mouseDownHandler);
})


function handleCalcCostTime() {
  const rectinfo = core.ImageRectInfo.createFrom({
    x: selectBoxRef.value.offsetLeft - previewImgRef.value.offsetLeft,
    y: selectBoxRef.value.offsetTop - previewImgRef.value.offsetTop,
    w: selectBoxRef.value.offsetWidth,
    h: selectBoxRef.value.offsetHeight,
    preview_width: previewImgRef.value.offsetWidth,
    preview_height: previewImgRef.value.offsetHeight,
    source_width: props.data.width,
    source_height: props.data.height,
  })

  StartAnalyse(rectinfo, threshold)
  NProgress.start()
}
// cropBtn.addEventListener('click', () => {
//   const sX = previewImgRef.value.offsetLeft - previewImgRef.value.offsetLeft;  // 区域选择框左侧位置
//   const sY = previewImgRef.value.offsetTop - previewImgRef.value.offsetTop;  // 区域选择框上方位置
//   const sW = previewImgRef.value.offsetWidth;  // 区域选择框宽度
//   const sH = previewImgRef.value.offsetHeight;  // 区域选择框高度
//   // 把图片截取到 canvas
//   canvasEl.getContext('2d').drawImage(previewImgRef.value, sX, sY, sW, sH , 0, 0, canvasEl.width, canvasEl.height);
//   // 把裁剪后的 canvas 图像转为 Blob
//   canvasEl.toBlob(blob => {
//     if (blob === null) return false;
//     imgFile = blob;
//   }, 'image/jpeg');
// });

function handleImageLoadSuccess() {
  console.log("handleImageLoadSuccess")
}


const mouseDownHandler = function (e:any) {
    e.stopPropagation()
    // Get the current mouse position
    console.log("mouseDownHandler")
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

const mouseMoveHandler = function (e: any) {
  console.log("mouseMoveHandler")
    // How far the mouse has been moved
    const dx = e.clientX - location.x;
    const dy = e.clientY - location.y;
    console.log(dx, dy)
    // Adjust the dimension of element
    // selectBoxRef.value.style.width = `${location.w + dx}px`;
    // selectBoxRef.value.style.height = `${location.h + dy}px`;
    selectBoxStyle.width = `${location.w + dx}px`;
    selectBoxStyle.height = `${location.h + dy}px`;
};

const mouseUpHandler = function () {
    // Remove the handlers of `mousemove` and `mouseup`
    document.removeEventListener('mousemove', mouseMoveHandler);
    document.removeEventListener('mouseup', mouseUpHandler);
};


onMounted(()=>{
  EventsOn("latency:analyse_start", ()=>{
    ElNotification({
      title: '进度提示',
      type: 'info',
      message: "开始数据分析",
    })
  })
  EventsOn("latency:analyse_filish", (res: Array<number>)=>{
    // processStatus.value = 0
    // rightContentRef.value.loadResponseTimeData(res)
    delayTimes.value = res
    ElNotification({
      title: '进度提示: 3/3',
      type: 'success',
      message: "数据处理完成",
    })
    NProgress.done()
  })

})

onUnmounted(()=>{
  EventsOff("latency:analyse_start")
  EventsOff("latency:analyse_filish")
})

/**
 * 判断是否竖屏
//  */
// function isVerticalScreen() {
//   return props.data.width < props.data.height
// }
// watch(props.data, value => {

// })

/**
 * 判断是否竖屏
 */
// const isVerticalScreen = computed(()=>{
//   return props.data.width < props.data.height
// })


function enableCalcButton() {
  calcButtonDisable.value = false
}

function setCalcButtonDisable(value: boolean) {
  calcButtonDisable.value = value
}

function setImagePlaceHolder() {
  previewImgRef.value.src = ""
}

function loadNewImage(info: core.ImageInfo) {
  console.log(previewImgRef.value)
  previewImgRef.value.src = info.path
  imageInfo.path = info.path
  imageInfo.width = info.height
  imageInfo.height = info.height
}

defineExpose({
  enableCalcButton,
  setCalcButtonDisable,
  setImagePlaceHolder,
  loadNewImage
})



</script>

<template>
  <div>
    <el-scrollbar height="calc(95vh)">
    <el-row justify="center" class="preview-content">
      <el-col :span="22">
        <span>标识检测区域</span>
        <div class="out-img-box">
          <img ref="previewImgRef" class="preview-img" draggable="false" :src="imageInfo.path == '' ? ' ./static/images/placeholder.png' : imageInfo.path" alt=""/>
          <div ref="selectBoxRef" :style="selectBoxStyle" class="s-move-content-header" id="select-box">
            <div ref="resizeTopRef" class="resizer resizer-t"></div>
            <div ref="resizeRightRef" class="resizer resizer-r"></div>
            <div ref="resizeBottomRef" class="resizer resizer-b"></div>
            <div ref="resizeLeftRef" class="resizer resizer-l"></div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row justify="end">
      <el-col :span="4">
        <el-button type="primary" @click="handleCalcCostTime" :disabled="calcButtonDisable" style="float:right">计算延迟</el-button>
      </el-col>
    </el-row>
    <el-row>
      <span class="data-item">检测到操作总数: {{ delayTimes.length }}</span>
      <!-- <span class="data-item">延迟: {{ delayTimes.join(", ") }}</span> -->
    </el-row>
    <el-row>
      <span class="data-item">计算延迟数据为: {{ delayTimes.join(", ") }}</span>
    </el-row>
  </el-scrollbar>
  </div>
</template>

<style>

.out-img-box {
    width: calc(70vw);
    height: calc(70vh);
    /* line-height: 600px; */
    text-align: center;
}

.preview-img {
    max-width: 100%;
    max-height: 100%;
    vertical-align: middle;
    
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

.preview-content {

}

</style>