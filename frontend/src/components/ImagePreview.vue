<script setup lang="ts">
import {reactive, ref, inject, Ref, onMounted, computed} from 'vue'

const selectBoxRef = ref()
const previewImgRef = ref()
const canvasRef = ref()
const data = [112,240, 256, 240, 224 ,240]
const url =
    'https://fuss10.elemecdn.com/e/5d/4a731a90594a4af544c0c25941171jpeg.jpeg'


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


function move(e:any) {
  // document.getElementById("local").innerHTML = e.pageY + '，' + e.pageX
  selectArea.mx = e.pageX;
  selectArea.my = e.pageY;
  paintReact()
}

function down(e:any) {
  // document.getElementById("down").innerHTML = e.pageY + '，' + e.pageX
  selectArea.ax = e.pageX;
  selectArea.ay = e.pageY;
  selectArea.paint = true
}

function up(e:any) {
  // document.getElementById("up").innerHTML = e.pageY + '，' + e.pageX
  selectArea.bx = e.pageX;
  selectArea.by = e.pageY;
  selectArea.paint = false
}

function paintReact() {
  const ctx = canvasRef.value.getContext("2d")
  ctx.clearRect(0, 0, 1000, 800)
  selectArea.width = selectArea.mx - selectArea.ax
  selectArea.height = selectArea.my - selectArea.ay
  console.log(selectArea.width, selectArea.height);
  if (selectArea.paint) {
      ctx.strokeRect(selectArea.ax, selectArea.ay, selectArea.width, selectArea.height)
  }
}

onMounted(()=>{
  previewImgRef.value.addEventListener('load', ()=>{
    selectBoxRef.value.style.display = 'block'
    selectBoxRef.value.style.top = previewImgRef.value.offsetTop + 'px'
    selectBoxRef.value.style.left = previewImgRef.value.offsetLeft + 'px'

    console.log(selectBoxRef.value.offsetHeight)
    console.log(selectBoxRef.value.style.top)
  })

  selectBoxRef.value.addEventListener('mousedown', (ev:any) => {
    console.log(ev)
    const X = ev.clientX - ev.target.offsetLeft;
    const Y = ev.clientY - ev.target.offsetTop;

    // 鼠标移动
    document.onmousemove = (ev:any) => {
      selectBoxRef.value.style.left = ev.clientX - X + 'px';
      selectBoxRef.value.style.top = ev.clientY - Y + 'px';

      // 限制选择框的拖动范围，禁止拖出图片区域
      // if (selectBoxRef.value.offsetLeft <= previewImgRef.value.offsetLeft) {
      //   selectBoxRef.value.style.left = previewImgRef.value.offsetLeft + 'px';
      // }
      // if (selectBoxRef.value.offsetLeft >= previewImgRef.value.offsetWidth - selectBoxRef.value.offsetWidth) {
      //   selectBoxRef.value.style.left = previewImgRef.value.offsetWidth - selectBoxRef.value.offsetWidth + 'px';
      // }
      // if (selectBoxRef.value.offsetTop <= previewImgRef.value.offsetTop) {
      //   selectBoxRef.value.style.top = previewImgRef.value.offsetTop + 'px';
      // }
      // if (selectBoxRef.value.offsetTop >= previewImgRef.value.offsetHeight - selectBoxRef.value.offsetHeight) {
      //   selectBoxRef.value.style.top = previewImgRef.value.offsetHeight - selectBoxRef.value.offsetHeight + 'px';
      // }
    }

    document.onmouseup = (ev:any) => {
      console.log(ev)
      console.log(selectBoxRef.value.offsetHeight)
      console.log(selectBoxRef.value.style.top)
      document.onmousemove = null;
    }
    return false;
  })

  // canvasRef.value.addEventListener('mousemove', move, false)
  // canvasRef.value.addEventListener('mousedown', down, false)
  // canvasRef.value.addEventListener('mouseup', up, false)
})


function getImage() {
  const sX = selectBoxRef.value.offsetLeft - previewImgRef.value.offsetLeft;  // 区域选择框左侧位置
  const sY = selectBoxRef.value.offsetTop - previewImgRef.value.offsetTop;  // 区域选择框上方位置
  const sW = selectBoxRef.value.offsetWidth;  // 区域选择框宽度
  const sH = selectBoxRef.value.offsetHeight;  // 区域选择框高度
  const pW = previewImgRef.value.offsetWidth;  // 区域选择框宽度
  const pH = previewImgRef.value.offsetHeight;  // 区域选择框高度
  console.log(sX, sY, sW, sH, pW, pH)

  // const destArea = 0
  // 把图片截取到 canvas
  // canvasEl.getContext('2d').drawImage(previewImgRef.value, sX, sY, sW, sH , 0, 0, canvasEl.width, canvasEl.height);
  // // 把裁剪后的 canvas 图像转为 Blob
  // canvasEl.toBlob(blob => {
  //   if (blob === null) return false;
  //   imgFile = blob;
  // }, 'image/jpeg');
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
</script>

<template>
    <div>
    <!-- <el-empty description="description" /> -->
    <el-button @click="getImage">getImageXY</el-button>
    
    <div class="out-img-bobx">
        <!-- <span class="demonstration">111</span> -->
        <img ref="previewImgRef" class="preview-img" draggable="false" src="../assets/images/0001.png" alt=""/>
        <div ref="selectBoxRef" id="select-box">
      </div>
        <!-- <canvas ref="canvasRef" id="myCanvas" width="1000" height="800" style="border:1px solid #d3d3d3;">
			Your browser does not support the HTML5 canvas tag.
		</canvas> -->
		<!-- <span>local</span>
		<p id="local"></p>
		<span>down</span>
		<p id="down"></p>
		<span>up</span>
		<p id="up"></p> -->


    </div>
        

    </div>

    
</template>

<style scoped>

.out-img-bobx {
    width: 600px;
    height: 400px;
    line-height: 400px;
    text-align: center;
}

.preview-img {
    max-width: 100%;
    max-height: 100%;
    vertical-align: middle;
    
}

#select-box {
  width: 200px;
  height: 200px;
  background: rgba(255, 255, 0, 0.4);
  position: absolute;
  display: none;
  cursor: move;
}

</style>