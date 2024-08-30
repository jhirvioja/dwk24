# dwk24

Exercises for University of Helsinki DevOps with Kubernetes 2024

Learn more @ [devopswithkubernetes.com](https://devopswithkubernetes.com/)

## Exercise 3.06 DBaaS vs DIY

Many companies opt for a DBaaS solution, primarily because of security reasons. In the early stages of database creation, misconfigurations can lead to vulnerabilities that malicious actors might exploit. I once in a hobby project forgot to disable public IP access for a database, and by the next morning someone had managed to start a crypto mining operation using said database. SQL is technically Turing complete, after all.

Human error is something that can happen to even the best of us. This needs to be taken in to account. In the cloud, organization policies and governance are typically in place to prevent such mistakes. For example, a simple yet effective policy to combat the aforementioned would be that a created managed database can't have a public IP address.

Many people take managed databases for granted. Once you've manually configured a database yourself, you start to appreciate the time savings that managed solutions offer. In this field time is money, and whether you use 80-120€ per hour to hire an external professional or employ your own, you will save at least a couple days of initial setup and config work by just simply creating a managed database.

Over time, you will save costs related to upkeep and updates. Any code you write today will eventually become legacy code. Even if you follow best practices and use Infrastructure as Code (IaC), automating backups with scripts and so on, these setups can fail when you least expect it. For those times, some sort of documentation with failover steps would be preferred. It is nice to have a managed database, where they take care of most of tasks related to this.

If whatever you are doing with your database provides successful, scale is one thing to keep an eye on. Without proper monitoring, issues like running out of disk space or failed backup scripts can go unnoticed. Cloud providers these days provide easy monitoring and alerting solutions, which require little to none upkeep compared to a custom solution.

Other things managed database takes care of is providing a clearer path for version upgrades, and patching security vulnerabilities and other problems when they arise. Depending on the choice of image and setup, images can be corrupted and vulnerable.

In general, cloud-managed databases are often preferred to cut costs, but there are always some cases where you want and need to have full control of the database. For example, if your data needs to have pretty much 100% availability at all times. Cloud managed infrastructure [can also have its mishaps](https://cloud.google.com/blog/products/infrastructure/details-of-google-cloud-gcve-incident). For hobby projects, DIY database management can be a good learning experience, especially if there are no major security concerns.
