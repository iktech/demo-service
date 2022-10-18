version = "1.0.0"
projectName = "demo-server"
repository = "iktech/demo-service"
image = ""

podTemplate(label: 'gopod', cloud: 'kubernetes', serviceAccount: 'jenkins',
  containers: [
    containerTemplate(name: 'docker', image: 'docker:dind', ttyEnabled: true, command: 'cat', privileged: true,
        envVars: [secretEnvVar(key: 'DOCKER_USERNAME', secretName: 'ikolomiyets-docker-hub-credentials', secretKey: 'username'),
    ]),
    containerTemplate(name: 'kubectl', image: 'roffe/kubectl', ttyEnabled: true, command: 'cat'),
  ],
  volumes: [
    secretVolume(secretName: 'ssh-home', mountPath: '/etc/.ssh'),
    secretVolume(secretName: 'ikolomiyets-docker-hub-credentials', mountPath: '/etc/.secret'),
    hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock')
  ]) {
    node('gopod') {
  		stage('Prepare') {
            checkout scm
            image="${repository}:${version}.${env.BUILD_NUMBER}"

            // Set up private key to access BitBucket
            sh "cat /etc/.ssh/id_rsa > ~/.ssh/id_rsa"
            sh "chmod 400 ~/.ssh/id_rsa"
  		}

		container('docker') {
			stage('Build Docker Image') {
				try {
					sh "docker build --build-arg VERSION=${version}.${env.BUILD_NUMBER} -t ${image} ."
					sh 'cat /etc/.secret/password | docker login --password-stdin --username $DOCKER_USERNAME'
					sh """
						docker push ${image}
						docker tag ${image} ${repository}:latest
						docker push ${repository}:latest
					"""
				} catch (error) {
					slackSend color: "danger", message: "Publishing Failed - ${env.JOB_NAME} build number ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
					throw error
				}
			}
        }

		stage('Publish the docker version to the Artifactor') {
			publishArtifact name: 'demo-service',
							 type: 'DockerImage',
							 stage: 'Development',
							 flow: 'Simple',
							 version: "${version}.${env.BUILD_NUMBER}"
		}

        stage('Tag Source Code') {
            checkout scm

            def repositoryCommitterEmail = "jenkins@iktech.io"
            def repositoryCommitterUsername = "jenkinsCI"
            values = version.tokenize(".")

            sh "git config user.email ${repositoryCommitterEmail}"
            sh "git config user.name '${repositoryCommitterUsername}'"
            sh "git tag -d v${values[0]} || true"
            sh "git push origin :refs/tags/v${values[0]}"
            sh "git tag -d v${values[0]}.${values[1]} || true"
            sh "git push origin :refs/tags/v${values[0]}.${values[1]}"
            sh "git tag -d v${version} || true"
            sh "git push origin :refs/tags/v${version}"

            sh "git tag -fa v${values[0]} -m \"passed CI\""
            sh "git tag -fa v${values[0]}.${values[1]} -m \"passed CI\""
            sh "git tag -fa v${version} -m \"passed CI\""
            sh "git tag -a v${version}.${env.BUILD_NUMBER} -m \"passed CI\""
            sh "git push -f --tags"

            milestone(5)
        }

		stage('Push the version from the Development to Integration Test stage') {
			pushArtifact name: 'demo-service', stage: 'Development'
		}

        stage('Deploy Latest') {
            container('kubectl') {
                try {
                    dir('deployment/uat') {
                        sh """
                            sed 's/_VERSION_/${tag}/g' namespace.yaml | kubectl apply -f -
                        """
                    }
                } catch (error) {
                    echo 'Deployment failed: ' + error.toString()
                    step([$class: 'Mailer',
                        notifyEveryUnstableBuild: true,
                        recipients: emailextrecipients([[$class: 'CulpritsRecipientProvider'],
                                                        [$class: 'DevelopersRecipientProvider']]),
                        sendToIndividuals: true])
                    slackSend color: "danger", message: "Build Failure - ${env.JOB_NAME} build number ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
                    throw error
                }
                milestone(6)
            }
        }
	}
}

properties([[
    $class: 'BuildDiscarderProperty',
    strategy: [
        $class: 'LogRotator',
        artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10']
    ]
]);
