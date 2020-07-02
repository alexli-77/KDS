

### KDS

> A self-adaptive and policy-based dynamic scheduler for applications in Kubernetes

#### Introduction

This project is based on Descheduler. The Descheduler aimed to carry out the research and development of load rebalancing based on cloud computing services. Service containerization has become well known as the mainstream deployment of distributed services. Kubernetes is the mainstream container arrangement tool, which provides good dynamic scaling, deployment, release and other functions for containerized services. However, until now, the load balancing of distributed services is still a part that needs to be continuously optimized.We need to work out a scheme that can guarantee both the dynamic scalability of the service and the dynamic load balancing of the service.

However, KDS is promoted on the basis of Descheduler,  and KDS considered and solved two problems:

1.Descheduler is used to solve the load balancing problem of stateless services. For stateful services, that is, those with local storage, the security of data needs to be considered during scheduling. 

2.we should ensure that the state is automatically adjusted and ultimately balanced to achieve the cluster load, rather than just providing the final state artificially.

This project consists of two parts, KDS and operator. This Github project is the KDS.

Notes：
1.The current load balancing solution provided by the Kubernetes platform is suitable for cluster operation. Since we periodically detect the cluster state and detect the number of pods being deployed during the migration to ensure that the cluster load is not affected, application services can be allocated to new physical resources to be added later in the run to address the load imbalance.

2.For service applications with storage capabilities, we can take the approach of multiple copies, keeping the same data on each copy and performing data copy operations in the initialization container to ensure that the data can be recovered after scheduling. The data here includes data in memory and on disk.

Research Platform
linux

Programming Language
go, Java, C++.

Research Topic
1 A self-adaptive and policy-based dynamic scheduler for applications in Kubernetes

Build and Run
Build descheduler:

$ make and run descheduler:

$ ./_output/bin/descheduler --kubeconfig --policy-config-file If you want more information about what descheduler is doing add -v 1 to the command line

For more information about available options run:

$ ./_output/bin/descheduler --help

