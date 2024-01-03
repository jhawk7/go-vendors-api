pipeline {
  agent {
    kubernetes {
      yaml '''
      apiVersion: v1
      kind: Pod
      metadata:
        name: kaniko
        namespace: devops
      spec:
        nodeSelector:
          kubernetes.io/hostname: node7
        containers:
        - name: kaniko
          image: gcr.io/kaniko-project/executor:debug
          imagePullPolicy: Always
          command:
          - /busybox/cat
          tty: true
          volumeMounts:
            - name: docker-config
              mountPath: /kaniko/.docker/
          envFrom:
          - secretRef:
              name: github-secret
        restartPolicy: Never
        volumes:
          - name: docker-config
            secret:
              secretName: docker-credentials
              items:
                - key: .dockerconfigjson
                  path: config.json
      '''
    }
  }

  stages {
    stage("build/deploy") {
      steps {
        checkout scm
        container(name: 'kaniko', shell: '/busybox/sh') {
          withEnv(['PATH+EXTRA=/busybox:/kaniko']) {
            sh '''#!/busybox/sh
            /kaniko/executor --dockerfile=Dockerfile \
            --context=`pwd` \
            --cache=true \
            --custom-platform=linux/arm64 \
            --destination=jhawk7/go-vendors-api:$BUILD_ID \
            --destination=jhawk7/go-vendors-api:latest'''
          }
        }   
      }
    }
  }
}