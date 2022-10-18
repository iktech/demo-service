version = "1.0.0"
projectName = "demo-server"
repository = "iktech/demo-service"
image = ""

podTemplate(label: 'gopod', cloud: 'kubernetes', serviceAccount: 'jenkins',
  containers: [
    containerTemplate(name: 'buildkit', image: 'moby/buildkit:master', ttyEnabled: true, privileged: true),
    containerTemplate(name: 'kubectl', image: 'roffe/kubectl', ttyEnabled: true, command: 'cat'),
  ],
  volumes: [
    secretVolume(secretName: 'ssh-home', mountPath: '/etc/.ssh'),
    secretVolume(secretName: 'docker-config', mountPath: '/etc/.secret')
  ]) {
    node('gopod') {
  		stage('Prepare') {
            checkout scm
            image="${repository}:${version}.${env.BUILD_NUMBER}"

            // Set up private key to access BitBucket
            sh "cat /etc/.ssh/id_rsa > ~/.ssh/id_rsa"
            sh "chmod 400 ~/.ssh/id_rsa"
  		}

		container('buildkit') {
			stage('Build Docker Image') {
				try {
					sh """
						mkdir -p /root/.docker
						cp /etc/.secret/config.json /root/.docker/

						wget https://raw.githubusercontent.com/anchore/grype/main/install.sh
						chmod 755 install.sh

						./install.sh
						mv bin/grype /usr/local/bin/

						sed 's/__VERSION__/${version}.${env.BUILD_NUMBER}/g' internal/adapters/left/version/handler.go > internal/adapters/left/version/handler.go.tmp
						mv internal/adapters/left/version/handler.go.tmp internal/adapters/left/version/handler.go

						buildctl build --frontend dockerfile.v0 --local context=. --local dockerfile=. --export-cache type=local,dest=/tmp/buildkit/cache --output type=oci,dest=/tmp/image.tar --opt network=host --opt build-arg:VERSION=${version}.${env.BUILD_NUMBER}
					"""
				} catch (error) {
					slackSend color: "danger", message: "Publishing Failed - ${env.JOB_NAME} build number ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
					throw error
				}
			}

			stage('Publish the docker version to the Artifactor') {
				publishArtifact name: 'demo-service',
								 type: 'DockerImage',
								 stage: 'Development',
								 flow: 'Simple',
								 version: "${version}.${env.BUILD_NUMBER}"
			}


			stage('Scan Docker Image') {
				try {
					sh """
						wget https://iktech-public-dl.s3.eu-west-1.amazonaws.com/grype/grype.tmpl
						GRYPE_MATCH_GOLANG_USING_CPES=false /usr/local/bin/grype oci-archive:/tmp/image.tar -f high --scope all-layers -o template --file report.html -t grype.tmpl
						ls -l
					"""
				} catch (error) {
					slackSend color: "danger", message: "Grype scan has failed - ${env.JOB_NAME} build number ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
					throw error
				} finally {
					publishHTML (target : [allowMissing: false,
                     alwaysLinkToLastBuild: true,
                     reportDir: '',
                     keepAll: true,
                     reportFiles: 'report.html',
                     reportName: 'Grype Scan Report',
                     reportTitles: 'Grype Scan Report'])
				}
				milestone(2)
			}
        }

		if (env.BRANCH_NAME == "master") {
			container('buildkit') {
				stage('Push Docker Images to the registry') {
					try {
						sh """
							buildctl build --frontend dockerfile.v0 --local context=. --local dockerfile=. --import-cache type=local,src=/tmp/buildkit/cache --output type=image,name=${image},push=true --opt network=host --opt build-arg:VERSION=${version}.${env.BUILD_NUMBER}
							buildctl build --frontend dockerfile.v0 --local context=. --local dockerfile=. --import-cache type=local,src=/tmp/buildkit/cache --output type=image,name=${repository}:latest,push=true --opt network=host --opt build-arg:VERSION=${version}.${env.BUILD_NUMBER}
						"""
					} catch (error) {
						slackSend color: "danger", message: "Pushing Docker Image has failed - ${env.JOB_NAME} build number ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
						throw error
					}
					milestone(3)
				}
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
						dir('deployment') {
							sh """
								sed 's/_VERSION_/${version}.${env.BUILD_NUMBER}/g' namespace.yaml | kubectl apply -f -
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
}

properties([[
    $class: 'BuildDiscarderProperty',
    strategy: [
        $class: 'LogRotator',
        artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10']
    ]
]);

