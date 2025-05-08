<template>
  <div class="home-container">
    <div class="hero-section">
      <div class="hero-content">
        <div class="hero-text">
          <h1>欢迎来到我的世界</h1>
          <p class="subtitle">探索、创造、分享</p>
        </div>

        <section class="features-section">
          <div class="feature-card">
            <h2>服务器列表</h2>
            <p>查看所有可用的服务器</p>
            <router-link to="/server_list" class="feature-link">
              立即查看
            </router-link>
          </div>

          <div class="feature-card">
            <h2>更多功能</h2>
            <p>即将推出更多精彩内容</p>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomePage'
}
</script>

<style scoped>
.home-container {
  min-height: 100%;
  background-color: #ffffff;
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow-x: hidden;
}

.dark-theme .home-container {
  background-color: #2c2c2c;
}

.hero-section {
  position: relative;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: left;
  overflow: hidden;
  background-image: url('@/assets/images/minecraft-bg.png');
  background-size: cover;
  background-repeat: no-repeat;
  width: 100%;
}

.hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0.5) 0%,
    rgba(0, 0, 0, 0.3) 50%,
    rgba(0, 0, 0, 0.5) 100%
  );
  z-index: 1;
}

.hero-content {
  position: relative;
  z-index: 2;
  width: 100%;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.hero-text {
  flex: 1;
  max-width: 500px;
}

.hero-section h1 {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #ffffff;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  white-space: nowrap;
}

.subtitle {
  font-size: 28px;
  color: #ffffff;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  margin-bottom: 40px;
}

.features-section {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  padding: 0;
  margin: 0;
  background-color: transparent;
  flex: 1;
  max-width: 600px;
}

.feature-card {
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 18px;
  padding: 30px;
  transition: transform 0.3s ease;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.dark-theme .feature-card {
  background-color: rgba(29, 29, 31, 0.9);
}

.feature-card:hover {
  transform: translateY(-5px);
}

.feature-card h2 {
  font-size: 24px;
  margin-bottom: 15px;
  color: #1d1d1f;
}

.dark-theme .feature-card h2 {
  color: #f5f5f7;
}

.feature-card p {
  color: #86868b;
  margin-bottom: 20px;
}

.dark-theme .feature-card p {
  color: #a1a1a6;
}

.feature-link {
  display: inline-block;
  padding: 12px 24px;
  background-color: #0071e3;
  color: white;
  text-decoration: none;
  border-radius: 20px;
  transition: background-color 0.3s ease;
}

.feature-link:hover {
  background-color: #0077ed;
}

@media (max-width: 768px) {
  .hero-section {
    height: 100vh;
  }

  .hero-content {
    flex-direction: column;
    text-align: center;
    gap: 40px;
  }

  .hero-text {
    max-width: 100%;
  }

  .features-section {
    max-width: 100%;
  }

  .hero-section h1 {
    font-size: 32px;
    white-space: nowrap;
  }

  .subtitle {
    font-size: 20px;
  }
}
</style>
