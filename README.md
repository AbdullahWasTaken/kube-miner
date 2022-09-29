# kube-miner

Something on the lines of building a complete package for learning on inherent structures found in the microservices architecture.
Functionally involving extarction of operational state of clusters running kubernetes, augmenting this using structural and accessibility patterns to gain insights over the general structure.

## Usage:
```
kube-miner collects the operational state of a k8s cluster.

Usage:
  ./kube-miner [flags]

Flags:
  -kubeconfig string      (optional) absolute path to the kubeconfig file (default "/home/<USER>/.kube/config")
  -outputPath string      (optional) relative path to where the collected data will be stored (default "state")

Example:
  ./kube-miner -kubeconfig <abs. path to kubeconfig> -outputPath <rel. path to output (need not exist)>
```
