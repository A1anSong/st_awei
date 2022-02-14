<template>
  <a-layout>
    <Header/>
    <a-layout :style="{ marginTop:'64px' }">
      <Sider/>
      <a-layout>
        <router-view/>
        <Footer/>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script>
import Header from "@/components/admin/Header";
import Sider from "@/components/admin/Sider";
import Footer from "@/components/admin/Footer";

export default {
  name: "Admin",
  components: {
    Header,
    Sider,
    Footer,
  },
}
</script>

<style scoped>

</style>