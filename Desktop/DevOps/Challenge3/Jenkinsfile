pipeline {
    agent any

    options {
        timestamps()
    }

    parameters {
        string(name: 'IMAGE_NAME', defaultValue: 'ttl.sh/abhi-challenge3:2h', description: 'Fully-qualified image name to build, push, and deploy.')
        string(name: 'DOCKER_VM_HOST', defaultValue: '', description: 'SSH host or IP address of the Docker VM.')
        string(name: 'DOCKER_VM_USER', defaultValue: 'ubuntu', description: 'SSH user on the Docker VM.')
        string(name: 'SSH_CREDENTIALS_ID', defaultValue: 'docker-vm-ssh-key', description: 'Jenkins SSH private key credential ID.')
    }

    environment {
        CONTAINER_NAME = 'challenge3'
        DOCKER_BUILDKIT = '1'
    }

    stages {
        stage('Build image') {
            steps {
                sh '''
                    set -eu
                    docker build --platform linux/amd64 -t "${IMAGE_NAME}" .
                '''
            }
        }

        stage('Push image') {
            steps {
                sh '''
                    set -eu
                    docker push "${IMAGE_NAME}"
                '''
            }
        }

        stage('Deploy on Docker VM') {
            steps {
                script {
                    if (!params.DOCKER_VM_HOST?.trim()) {
                        error('Set DOCKER_VM_HOST before running the deploy stage.')
                    }
                }

                sshagent(credentials: [params.SSH_CREDENTIALS_ID]) {
                    sh '''
                        set -eu
                        ssh -o StrictHostKeyChecking=no "${DOCKER_VM_USER}@${DOCKER_VM_HOST}" "
                            set -eu
                            docker pull '${IMAGE_NAME}'
                            docker rm -f '${CONTAINER_NAME}' >/dev/null 2>&1 || true
                            docker run -d --name '${CONTAINER_NAME}' --restart unless-stopped -p 4444:4444 '${IMAGE_NAME}'
                        "
                    '''
                }
            }
        }
    }
}
