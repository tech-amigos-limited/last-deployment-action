# get-deployment-action

Usage 
```
  steps:
    - name: Search deployments
    id: deployments
    uses: docker://techamigos/last-deployment-action:latest
    with: 
        github-token: ${{ secrets.GRID_GIT_TOKEN }}
        ref: ${{ github.head_ref }}
        repo: ${{ github.repository }}
```

Availble outputs
```
  ${{ steps.deployments.outputs.last_deployment_id }}
  ${{ steps.deployments.outputs.last_status }} 
```

