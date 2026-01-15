<template>
  <div class="fullwidth-page">
    <div class="card" style="max-width: 600px; margin: 40px auto;">
      <div class="page-header">
        <h2>{{ $t('import.title') }}</h2>
      </div>
      
      <p class="text-muted">{{ $t('import.description') }}</p>
      
      <form @submit.prevent="handleUpload">
        <div class="form-group">
          <label for="file_vtt">{{ $t('import.charVTT') }}</label>
          <input 
            type="file" 
            id="file_vtt" 
            class="form-control" 
            @change="onFileChange($event, 1)"
            accept=".json,.csv"
            required
          />
          <small class="form-text">{{ $t('import.vttHelp') }}</small>
        </div>
        
        <div class="form-group">
          <label for="file_csv">{{ $t('import.charCSV') }}</label>
          <input 
            type="file" 
            id="file_csv" 
            class="form-control" 
            @change="onFileChange($event, 2)"
            accept=".json,.csv"
          />
          <small class="form-text">{{ $t('import.csvHelp') }}</small>
        </div>
        
        <button 
          type="submit" 
          class="btn btn-primary" 
          style="width: 100%; margin-top: 10px;"
          :disabled="!file_vtt"
        >
          {{ $t('import.upload') }}
        </button>
      </form>
      
      <div v-if="error" class="badge badge-danger" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
        {{ error }}
      </div>
      
      <div v-if="success" class="badge badge-success" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
        {{ success }}
      </div>
    </div>
  </div>
</template>

<script>
import API from "../utils/api";

export default {
  data() {
    return {
      file_vtt: null,
      file_csv: null,
      error: "",
      success: "",
    };
  },
  computed: {
    hasInvalidFileType() {
      return !this.isValidFileType(this.file_vtt) || (this.file_csv && !this.isValidFileType(this.file_csv));
    },
  },
  methods: {
    onFileChange(event, fileNumber) {
      const file = event.target.files[0];
      if (fileNumber === 1) this.file_vtt = file;
      if (fileNumber === 2) this.file_csv = file;

      // Validate file type
      if (!this.isValidFileType(file)) {
        this.error = this.$t('import.invalidFileType');
        return;
      }
      this.error = ""; // Clear any previous error
    },
    isValidFileType(file) {
      if (!file) return false;
      const allowedTypes = ["application/json", "text/csv"];
      return allowedTypes.includes(file.type);
    },
    async handleUpload() {
      try {
        const formData = new FormData();
        formData.append("file_vtt", this.file_vtt);
        if (this.file_csv) formData.append("file_csv", this.file_csv);

        const token = localStorage.getItem("token"); // Get token from storage
        const response = await API.post("/api/importer/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
            "Authorization": `Bearer ${token}`, // Include token in the header
          },
        });

        this.success = this.$t('import.uploadSuccess');
        this.error = "";
      } catch (err) {
        this.error = err.response?.data?.error || this.$t('import.uploadFailed');
        this.success = "";
      }
    },
  },
};
</script>

<style scoped>
.text-muted {
  color: #6c757d;
  margin-bottom: 20px;
}

.form-text {
  display: block;
  margin-top: 0.25rem;
  color: #6c757d;
  font-size: 0.875rem;
}
</style>
