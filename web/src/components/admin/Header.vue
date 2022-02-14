<template>
  <a-layout-header :style="{ position: 'fixed', zIndex: 1, width: '100%' }">
    <div class="logo"/>
    <a-menu
        v-model:selectedKeys="selectedKeys"
        theme="dark"
        mode="horizontal"
        :style="{ lineHeight: '64px' }"
    >
      <a-menu-item key="Stripe">Stripe</a-menu-item>
      <a-menu-item key="EUR">EUR</a-menu-item>
    </a-menu>
  </a-layout-header>
</template>

<script>
export default {
  name: "Header",
  data() {
    return {
      selectedKeys: ['EUR'],
    }
  },
}
</script>

<style scoped>
.logo {
  width: 32px;
  height: 32px;
  background: url("~@/assets/logo.png");
  background-size: contain;
  margin: 16px 24px 16px 0;
  float: left;
}
</style>