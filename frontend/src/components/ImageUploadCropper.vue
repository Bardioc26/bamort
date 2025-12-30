<template>
  <div class="image-upload-container">
    <button @click="showDialog = true" class="btn-upload">
      {{ $t('character.uploadImage') }}
    </button>

    <div v-if="showDialog" class="modal-overlay" @click.self="closeDialog">
      <div class="modal-content image-upload-modal">
        <div class="modal-header">
          <h3>{{ $t('character.uploadImage') }}</h3>
          <button @click="closeDialog" class="close-button">&times;</button>
        </div>
        
        <div class="modal-body">
          <div v-if="!imageLoaded" class="upload-area">
            <input 
              type="file" 
              ref="fileInput" 
              @change="handleFileSelect" 
              accept="image/*"
              style="display: none"
            />
            <button @click="$refs.fileInput.click()" class="btn-select-file">
              {{ $t('character.selectImage') }}
            </button>
          </div>

          <div v-else class="crop-area">
            <div class="crop-controls">
              <label>
                <input type="radio" v-model="cropShape" value="rect" />
                {{ $t('character.cropRect') }}
              </label>
              <label>
                <input type="radio" v-model="cropShape" value="round" />
                {{ $t('character.cropRound') }}
              </label>
            </div>

            <div class="canvas-container">
              <canvas 
                ref="canvas" 
                @mousedown="startDrag"
                @mousemove="drag"
                @mouseup="endDrag"
                @mouseleave="endDrag"
              ></canvas>
            </div>

            <div class="preview-container">
              <h4>{{ $t('character.preview') }}</h4>
              <canvas ref="previewCanvas" width="400" height="400"></canvas>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeDialog" class="btn-cancel">{{ $t('cancel') }}</button>
          <button v-if="imageLoaded" @click="resetImage" class="btn-secondary">
            {{ $t('character.changeImage') }}
          </button>
          <button 
            v-if="imageLoaded" 
            @click="uploadImage" 
            class="btn-primary"
            :disabled="isUploading"
          >
            <span v-if="!isUploading">{{ $t('character.saveImage') }}</span>
            <span v-else>{{ $t('uploading') }}...</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
/* All common styles moved to main.css */

/* ImageUploadCropper specific styles */
.btn-upload {
  padding: 8px 16px;
  background-color: var(--primary-color);

  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-upload:hover {
  opacity: 0.9;
}

.image-upload-modal {
  min-width: 700px;
  max-width: 90vw;
}

.upload-area {
  text-align: center;
  padding: 40px;
}

.btn-select-file {
  padding: 12px 24px;
  background-color: var(--primary-color);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.crop-area {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.crop-controls {
  display: flex;
  gap: 20px;
  justify-content: center;
}

.crop-controls label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.canvas-container {
  display: flex;
  justify-content: center;
  border: 1px solid #ccc;
  background-color: #f5f5f5;
  overflow: auto;
  max-height: 500px;
}

.canvas-container canvas {
  cursor: crosshair;
}

.preview-container {
  text-align: center;
}

.preview-container h4 {
  margin-bottom: 10px;
}

.preview-container canvas {
  border: 1px solid #ccc;
  background-color: white;
}
</style>

<script>
import API from '../utils/api'

export default {
  name: 'ImageUploadCropper',
  props: {
    characterId: {
      type: [String, Number],
      required: true
    }
  },
  data() {
    return {
      showDialog: false,
      imageLoaded: false,
      isUploading: false,
      cropShape: 'rect',
      image: null,
      cropStart: null,
      cropEnd: null,
      isDragging: false,
      cropWidth: 400,
      cropHeight: 400
    }
  },
  methods: {
    handleFileSelect(event) {
      const file = event.target.files[0]
      if (!file) return

      const reader = new FileReader()
      reader.onload = (e) => {
        const img = new Image()
        img.onload = () => {
          this.image = img
          this.imageLoaded = true
          this.$nextTick(() => {
            this.initCanvas()
          })
        }
        img.src = e.target.result
      }
      reader.readAsDataURL(file)
    },
    
    initCanvas() {
      const canvas = this.$refs.canvas
      const ctx = canvas.getContext('2d')
      
      canvas.width = this.image.width
      canvas.height = this.image.height
      
      ctx.drawImage(this.image, 0, 0)
      
      // Initialize crop area in center
      this.cropStart = {
        x: Math.max(0, (this.image.width - this.cropWidth) / 2),
        y: Math.max(0, (this.image.height - this.cropHeight) / 2)
      }
      this.cropEnd = {
        x: this.cropStart.x + this.cropWidth,
        y: this.cropStart.y + this.cropHeight
      }
      
      this.drawCropOverlay()
      this.updatePreview()
    },
    
    drawCropOverlay() {
      const canvas = this.$refs.canvas
      const ctx = canvas.getContext('2d')
      
      // Redraw image
      ctx.drawImage(this.image, 0, 0)
      
      // Draw semi-transparent overlay
      ctx.fillStyle = 'rgba(0, 0, 0, 0.5)'
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      
      const width = this.cropEnd.x - this.cropStart.x
      const height = this.cropEnd.y - this.cropStart.y
      
      // Clear crop area
      ctx.clearRect(this.cropStart.x, this.cropStart.y, width, height)
      ctx.drawImage(
        this.image,
        this.cropStart.x, this.cropStart.y, width, height,
        this.cropStart.x, this.cropStart.y, width, height
      )
      
      // Draw border
      ctx.strokeStyle = '#00ff00'
      ctx.lineWidth = 2
      if (this.cropShape === 'round') {
        ctx.beginPath()
        const centerX = this.cropStart.x + width / 2
        const centerY = this.cropStart.y + height / 2
        const radius = Math.min(width, height) / 2
        ctx.arc(centerX, centerY, radius, 0, 2 * Math.PI)
        ctx.stroke()
      } else {
        ctx.strokeRect(this.cropStart.x, this.cropStart.y, width, height)
      }
    },
    
    updatePreview() {
      const previewCanvas = this.$refs.previewCanvas
      const ctx = previewCanvas.getContext('2d')
      
      const width = this.cropEnd.x - this.cropStart.x
      const height = this.cropEnd.y - this.cropStart.y
      
      ctx.clearRect(0, 0, 400, 400)
      
      if (this.cropShape === 'round') {
        ctx.save()
        ctx.beginPath()
        ctx.arc(200, 200, 200, 0, 2 * Math.PI)
        ctx.clip()
      }
      
      ctx.drawImage(
        this.image,
        this.cropStart.x, this.cropStart.y, width, height,
        0, 0, 400, 400
      )
      
      if (this.cropShape === 'round') {
        ctx.restore()
      }
    },
    
    startDrag(event) {
      const canvas = this.$refs.canvas
      const rect = canvas.getBoundingClientRect()
      
      // Calculate mouse position on the actual canvas (accounting for scaling)
      const scaleX = canvas.width / rect.width
      const scaleY = canvas.height / rect.height
      
      this.cropStart = {
        x: (event.clientX - rect.left) * scaleX,
        y: (event.clientY - rect.top) * scaleY
      }
      this.isDragging = true
    },
    
    drag(event) {
      if (!this.isDragging) return
      
      const canvas = this.$refs.canvas
      const rect = canvas.getBoundingClientRect()
      
      // Calculate mouse position on the actual canvas (accounting for scaling)
      const scaleX = canvas.width / rect.width
      const scaleY = canvas.height / rect.height
      
      this.cropEnd = {
        x: (event.clientX - rect.left) * scaleX,
        y: (event.clientY - rect.top) * scaleY
      }
      
      // Ensure minimum size
      const minSize = 50
      if (Math.abs(this.cropEnd.x - this.cropStart.x) < minSize ||
          Math.abs(this.cropEnd.y - this.cropStart.y) < minSize) {
        return
      }
      
      this.drawCropOverlay()
      this.updatePreview()
    },
    
    endDrag() {
      this.isDragging = false
    },
    
    async uploadImage() {
      this.isUploading = true
      
      try {
        // Get cropped image as base64
        const previewCanvas = this.$refs.previewCanvas
        const croppedImage = previewCanvas.toDataURL('image/png')
        
        await API.put(`/api/characters/${this.characterId}/image`, {
          image: croppedImage
        })
        
        this.$emit('image-updated', croppedImage)
        this.closeDialog()
        //alert(this.$t('character.imageUploadSuccess'))
      } catch (error) {
        console.error('Failed to upload image:', error)
        alert(this.$t('character.imageUploadError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isUploading = false
      }
    },
    
    resetImage() {
      this.imageLoaded = false
      this.image = null
      this.cropStart = null
      this.cropEnd = null
      if (this.$refs.fileInput) {
        this.$refs.fileInput.value = ''
      }
    },
    
    closeDialog() {
      this.showDialog = false
      this.resetImage()
    }
  },
  watch: {
    cropShape() {
      if (this.imageLoaded) {
        this.drawCropOverlay()
        this.updatePreview()
      }
    }
  }
}
</script>
