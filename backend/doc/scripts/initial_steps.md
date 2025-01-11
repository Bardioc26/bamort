# Install Dependencies
## Install the necessary packages:

cd backend
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get github.com/gin-contrib/cors
go get -u github.com/stretchr/testify
go get -u github.com/stretchr/testify/assert
go get -u gorm.io/driver/sqlite


cd /frontend
npm init vue@latest bamort
npm install
npm install axios vue-router@4
npm install vue-i18n@9
npm install pinia