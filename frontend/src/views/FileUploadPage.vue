<template>
  <div class="upload-page">
    <h2>Import Data</h2>
    <br/>
    <form @submit.prevent="handleUpload">
      <div>
        <label for="file_vtt">Char VTT:</label>
        <input type="file" id="file_vtt" @change="onFileChange($event, 1)" />
      </div>
      <div>
        <label for="file_csv">Char CSV (optional):</label>
        <input type="file" id="file_csv" @change="onFileChange($event, 2)" />
      </div>
      <button type="submit" :disabled="!file_vtt">Upload</button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="success" class="success">{{ success }}</p>
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
        this.error = "Invalid file type. Only .csv and .json files are allowed.";
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
        const response = await API.post("/api/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
            "Authorization": `Bearer ${token}`, // Include token in the header
          },
        });

        this.success = "Files uploaded successfully!";
        this.error = "";
      } catch (err) {
        this.error = err.response?.data?.error || "Failed to upload files.";
        this.success = "";
      }
    },
  },
};
</script>

<style>
.upload-page {
  padding: var(--padding-lg);
  margin-top: 2%;
}

.error {
  color: red;
}

.success {
  color: green;
}
</style>
