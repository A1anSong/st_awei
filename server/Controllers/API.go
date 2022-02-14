package Controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72/client"
	"goST/Models"
	"gorm.io/gorm/clause"
	"net/http"
	"sync"
	"time"
)

var taskMutex sync.RWMutex
var taskProcessing = false
var completeAccounts = 0
var completeBalances = 0
var completeInvoices = 0
var completePayouts = 0
var totalAccounts = 0
var totalBalances = 0
var totalInvoices = 0
var totalPayouts = 0
var updateAccountsStatusChannel = make(chan int)
var updateBalancesStatusChannel = make(chan int)
var updateInvoicesStatusChannel = make(chan int)
var updatePayoutsStatusChannel = make(chan int)

func GetDashAccounts(c *gin.Context) {
	var active int64
	var total int64
	StripeDB.Model(&Models.Account{}).Count(&active)
	StripeDB.Unscoped().Model(&Models.Account{}).Count(&total)
	c.JSON(http.StatusOK, gin.H{
		"active": active,
		"total":  total,
	})
}

func GetDashBalances(c *gin.Context) {
	var available int
	var pending int
	StripeDB.Table("balances").Select("sum(available)").Where("available > 0 and deleted_at is null").Row().Scan(&available)
	StripeDB.Table("balances").Select("sum(pending)").Where("pending > 0 and deleted_at is null").Row().Scan(&pending)
	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"pending":   pending,
	})
}

func GetDashInvoices(c *gin.Context) {
	var success int
	var inprocess int
	var failed int
	StripeDB.Raw(`select count(*)
from accounts
where deleted_at is null
  and id in (select account_id from invoices where status = 'paid');
`).Row().Scan(&success)
	StripeDB.Raw(`select count(*)
from accounts
where deleted_at is null
  and id not in (select distinct account_id from invoices where status = 'paid')
  and id not in (select distinct account_id from invoices where status = 'open');
`).Row().Scan(&inprocess)
	StripeDB.Raw(`select count(*)
from accounts
where deleted_at is null
  and id not in (select distinct account_id from invoices where status = 'paid')
  and id in (select distinct account_id from invoices where status = 'open');
`).Row().Scan(&failed)
	c.JSON(http.StatusOK, gin.H{
		"success":   success,
		"inprocess": inprocess,
		"failed":    failed,
	})
}

func GetDashPayouts(c *gin.Context) {
	var intransit int
	var intransitCount int
	var paid int
	var paidCount int
	StripeDB.Table("payouts").Select("sum(amount) , count(*)").Where("status = 'in_transit' and amount > 0").Row().Scan(&intransit, &intransitCount)
	StripeDB.Table("payouts").Select("sum(amount) , count(*)").Where("status = 'paid' and amount > 0").Row().Scan(&paid, &paidCount)
	c.JSON(http.StatusOK, gin.H{
		"intransit":      intransit,
		"intransitCount": intransitCount,
		"paid":           paid,
		"paidCount":      paidCount,
	})
}

func getTaskProcessing() bool {
	taskMutex.Lock()
	defer taskMutex.Unlock()
	return taskProcessing
}

func setTaskProcessing(t bool) {
	taskMutex.Lock()
	taskProcessing = t
	taskMutex.Unlock()
}

func updateAccountsInfo() {
	var accounts []Models.Account
	StripeDB.Find(&accounts)
	totalAccounts = len(accounts)
	for index, account := range accounts {
		sc := &client.API{}
		sc.Init(account.SecretKey, nil)
		stripeAccount, _ := sc.Account.Get()
		account.Acc = stripeAccount.ID
		account.ChargesEnabled = stripeAccount.ChargesEnabled
		account.PayoutsEnabled = stripeAccount.PayoutsEnabled
		StripeDB.Save(&account)
		if !account.ChargesEnabled || !account.PayoutsEnabled {
			StripeDB.Select(clause.Associations).Delete(&account)
		}
		completeAccounts = index + 1
	}
	updateAccountsStatusChannel <- 1
}

func updateBalancesInfo() {
	var accounts []Models.Account
	StripeDB.Find(&accounts)
	totalBalances = len(accounts)
	for index, account := range accounts {
		sc := &client.API{}
		sc.Init(account.SecretKey, nil)
		stripeBalance, _ := sc.Balance.Get(nil)
		var balance Models.Balance
		StripeDB.FirstOrInit(&balance, Models.Balance{AccountID: account.ID})
		balance.Available = int(stripeBalance.Available[0].Value)
		balance.Pending = int(stripeBalance.Pending[0].Value)
		StripeDB.Save(&balance)
		completeBalances = index + 1
	}
	updateBalancesStatusChannel <- 1
}

func updateInvoicesInfo() {
	var accounts []Models.Account
	StripeDB.Find(&accounts)
	totalInvoices = len(accounts)
	for index, account := range accounts {
		sc := &client.API{}
		sc.Init(account.SecretKey, nil)
		stripeInvoices := sc.Invoices.List(nil)
		for stripeInvoices.Next() {
			stripeInvoice := stripeInvoices.Invoice()
			var invoice Models.Invoice
			StripeDB.FirstOrInit(&invoice, Models.Invoice{In: stripeInvoice.ID})
			invoice.AutoAdvance = stripeInvoice.AutoAdvance
			if stripeInvoice.CollectionMethod != nil {
				invoice.CollectionMethod = string(*stripeInvoice.CollectionMethod)
			} else {
				invoice.CollectionMethod = ""
			}
			invoice.Currency = string(stripeInvoice.Currency)
			invoice.Description = stripeInvoice.Description
			invoice.HostedInvoiceURL = stripeInvoice.HostedInvoiceURL
			invoice.PeriodEnd = time.Unix(stripeInvoice.PeriodEnd, 0)
			invoice.PeriodStart = time.Unix(stripeInvoice.PeriodStart, 0)
			invoice.Status = string(stripeInvoice.Status)
			invoice.Total = int(stripeInvoice.Total)
			invoice.AccountID = account.ID
			StripeDB.Save(&invoice)
		}
		completeInvoices = index + 1
	}
	updateInvoicesStatusChannel <- 1
}

func updatePayoutsInfo() {
	var accounts []Models.Account
	StripeDB.Unscoped().Find(&accounts)
	totalPayouts = len(accounts)
	for index, account := range accounts {
		sc := &client.API{}
		sc.Init(account.SecretKey, nil)
		stripePayouts := sc.Payouts.List(nil)
		for stripePayouts.Next() {
			stripePayout := stripePayouts.Payout()
			var payout Models.Payout
			StripeDB.Unscoped().FirstOrInit(&payout, Models.Payout{Po: stripePayout.ID})
			payout.Amount = int(stripePayout.Amount)
			payout.ArrivalDate = time.Unix(stripePayout.ArrivalDate, 0)
			payout.Currency = string(stripePayout.Currency)
			if stripePayout.Description != nil {
				payout.Description = *stripePayout.Description
			} else {
				payout.Description = ""
			}
			payout.StatementDescriptor = stripePayout.StatementDescriptor
			payout.Status = string(stripePayout.Status)
			payout.AccountID = account.ID
			StripeDB.Save(&payout)
		}
		completePayouts = index + 1
	}
	updatePayoutsStatusChannel <- 1
}

func updateInfos() {
	go updateAccountsInfo()
	go updatePayoutsInfo()
	<-updateAccountsStatusChannel
	go updateBalancesInfo()
	go updateInvoicesInfo()
	<-updateBalancesStatusChannel
	<-updateInvoicesStatusChannel
	<-updatePayoutsStatusChannel
	time.Sleep(time.Minute * 5)
	completeAccounts = 0
	completeBalances = 0
	completeInvoices = 0
	completePayouts = 0
	totalAccounts = 0
	totalBalances = 0
	totalInvoices = 0
	totalPayouts = 0
	setTaskProcessing(false)
}

func GetUpdateInfos(c *gin.Context) {
	if getTaskProcessing() {
		c.JSON(http.StatusOK, gin.H{
			"taskProcessing":   true,
			"completeAccounts": completeAccounts,
			"completeBalances": completeBalances,
			"completeInvoices": completeInvoices,
			"completePayouts":  completePayouts,
			"totalAccounts":    totalAccounts,
			"totalBalances":    totalBalances,
			"totalInvoices":    totalInvoices,
			"totalPayouts":     totalPayouts,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"taskProcessing": false,
		})
	}
}

func PatchUpdateInfos(c *gin.Context) {
	if getTaskProcessing() {
		//TODO: 加入上次执行时间判断，返回501
		c.String(http.StatusPreconditionFailed, "已经有在进行的更新，将会获取更新状态")
	} else {
		setTaskProcessing(true)
		go updateInfos()
		c.String(http.StatusOK, "成功提交信息更新任务")
	}
}
