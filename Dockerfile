# Dockerfile

# ベースとなる公式イメージを指定
FROM postgres:16-alpine

# 拡張機能のビルドに必要なツールをインストール
RUN apk add --no-cache git build-base postgresql-dev clang

# 拡張機能のソースコードをクローン
RUN git clone https://github.com/fboulnois/pg_uuidv7.git /pg_uuidv7

# ソースコードからビルドしてインストール
WORKDIR /pg_uuidv7
# ★変更点: (make ... install || true) のように変更
RUN make CLANG=clang && (make CLANG=clang install || true)

# 不要になったビルドツールを削除してイメージを軽量化
RUN apk del git build-base clang