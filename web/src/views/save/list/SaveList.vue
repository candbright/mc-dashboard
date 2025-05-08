<template>
  <div class="save-list-container">
    <div class="page-header">
      <h1>存档列表</h1>
      <p class="subtitle">管理你的游戏存档和备份</p>
    </div>

    <div class="save-list-content">
      <div class="save-list-header">
        <div class="search-box">
          <input
            type="text"
            placeholder="搜索存档..."
            v-model="searchQuery"
            class="search-input"
          >
        </div>
        <button class="add-button">
          <span class="button-icon">+</span>
          添加存档
        </button>
      </div>

      <div class="save-list">
        <div class="save-item" v-for="(save, index) in filteredSaves" :key="index">
          <div class="save-info">
            <h3>{{ save.name }}</h3>
            <p class="save-meta">
              <span>最后修改: {{ save.lastModified }}</span>
              <span>大小: {{ save.size }}</span>
            </p>
          </div>
          <div class="save-actions">
            <button class="action-button">备份</button>
            <button class="action-button">恢复</button>
            <button class="action-button delete">删除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SaveList',

  data() {
    return {
      searchQuery: '',
      saves: [
        {
          name: '生存模式存档',
          lastModified: '2024-03-20 15:30',
          size: '256MB'
        },
        {
          name: '创造模式存档',
          lastModified: '2024-03-19 10:15',
          size: '128MB'
        }
      ]
    }
  },

  computed: {
    filteredSaves() {
      return this.saves.filter(save =>
        save.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      )
    }
  }
}
</script>

<style scoped>
.save-list-container {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 40px;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 600;
  color: #1d1d1f;
  margin-bottom: 10px;
}

.dark-theme .page-header h1 {
  color: #f5f5f7;
}

.subtitle {
  font-size: 18px;
  color: #86868b;
}

.dark-theme .subtitle {
  color: #a1a1a6;
}

.save-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.search-box {
  flex: 1;
  max-width: 400px;
}

.search-input {
  width: 100%;
  padding: 12px 20px;
  border: none;
  border-radius: 8px;
  background-color: #f5f5f7;
  font-size: 16px;
  transition: all 0.3s ease;
}

.dark-theme .search-input {
  background-color: #1d1d1f;
  color: #f5f5f7;
}

.search-input:focus {
  outline: none;
  box-shadow: 0 0 0 2px #0071e3;
}

.add-button {
  display: flex;
  align-items: center;
  padding: 12px 24px;
  background-color: #0071e3;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.add-button:hover {
  background-color: #0077ed;
}

.button-icon {
  margin-right: 8px;
  font-size: 20px;
}

.save-list {
  display: grid;
  gap: 20px;
}

.save-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background-color: #f5f5f7;
  border-radius: 12px;
  transition: transform 0.3s ease;
}

.dark-theme .save-item {
  background-color: #1d1d1f;
}

.save-item:hover {
  transform: translateY(-2px);
}

.save-info h3 {
  font-size: 18px;
  font-weight: 500;
  color: #1d1d1f;
  margin-bottom: 8px;
}

.dark-theme .save-info h3 {
  color: #f5f5f7;
}

.save-meta {
  display: flex;
  gap: 20px;
  color: #86868b;
  font-size: 14px;
}

.dark-theme .save-meta {
  color: #a1a1a6;
}

.save-actions {
  display: flex;
  gap: 10px;
}

.action-button {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  background-color: #e8e8ed;
  color: #1d1d1f;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.dark-theme .action-button {
  background-color: #2c2c2c;
  color: #f5f5f7;
}

.action-button:hover {
  background-color: #d2d2d7;
}

.dark-theme .action-button:hover {
  background-color: #3c3c3c;
}

.action-button.delete {
  background-color: #ff3b30;
  color: white;
}

.action-button.delete:hover {
  background-color: #ff453a;
}

@media (max-width: 768px) {
  .save-list-header {
    flex-direction: column;
    gap: 20px;
  }

  .search-box {
    max-width: 100%;
  }

  .save-item {
    flex-direction: column;
    gap: 20px;
  }

  .save-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
