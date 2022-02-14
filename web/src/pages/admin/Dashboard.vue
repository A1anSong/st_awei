<template>
  <a-layout-content :style="{ padding: '0 50px', marginTop: '64px' }">
    <div :style="{ background: '#fff', padding: '24px', minHeight: '380px' }">
      <a-typography-title :level="4">总览</a-typography-title>
      <a-row>
        <AccountsInfo/>
        <BalancesInfo/>
      </a-row>
      <a-row>
        <InvoicesInfo/>
        <PayoutsInfo/>
      </a-row>
      <a-divider/>
      <a-typography-title :level="4">信息数据</a-typography-title>
      <UpdateAllInfo/>
    </div>
  </a-layout-content>
</template>

<script>
import AccountsInfo from "@/components/admin/dashboard/AccountsInfo";
import BalancesInfo from "@/components/admin/dashboard/BalancesInfo";
import InvoicesInfo from "@/components/admin/dashboard/InvoicesInfo";
import PayoutsInfo from "@/components/admin/dashboard/PayoutsInfo";
import UpdateAllInfo from "@/components/admin/dashboard/UpdateAllInfo";

export default {
  name: "Dashboard",
  components: {
    AccountsInfo,
    BalancesInfo,
    InvoicesInfo,
    PayoutsInfo,
    UpdateAllInfo,
  },
}
</script>

<style scoped>

</style>