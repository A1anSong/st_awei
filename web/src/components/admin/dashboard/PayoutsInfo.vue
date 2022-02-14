<template>
  <a-col :span="12">
    <a-spin :spinning="loading">
      <a-statistic title="出款：In transit/Paid" :value="payoutsInfoText"/>
    </a-spin>
  </a-col>
</template>

<script>
import {notification} from "ant-design-vue";

const axios = require('axios').default

export default {
  name: "PayoutsInfo",
  data() {
    return {
      loading: true,
      intransit: null,
      intransitCount: null,
      paid: null,
      paidCount: null
    }
  },
  computed: {
    payoutsInfoText() {
      return '€' + this.intransit + '(' + this.intransitCount + ')/€' + this.paid + '(' + this.paidCount + ')'
    },
  },
  methods: {
    getPayoutsInfo() {
      axios.get('/api/dashPayouts')
          .then(response => {
            this.intransit = response.data.intransit / 100
            this.intransitCount = response.data.intransitCount
            this.paid = response.data.paid / 100
            this.paidCount = response.data.paidCount
            this.loading = false
          })
          .catch(error => {
            notification.error({
              message: '获取出款信息错误',
              description: error.message,
              duration: 0,
            })
          })
    },
  },
  mounted() {
    this.getPayoutsInfo()
  },
}
</script>

<style scoped>

</style>