/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "notify",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new sst.aws.Vpc("MyVpc", { nat: "managed" });
    const database = new sst.aws.Postgres("MyDatabase", {
      vpc,
    });
    const cluster = new sst.aws.Cluster("MyCluster", {
      vpc,
    });

    cluster.addService("NotifyServer", {
      link: [database],
      cpu: "1 vCPU",
      memory: "2 GB",
      image: {
        context: "./", // DockerFile Path
      },
      public: {
        ports: [{ listen: "80/http", forward: "8080/http" }],
      },
      dev: {
        command: "go run ./cmd/main.go",
      },
      environment: {
        NOTIFY_DB_URL: "",
        NOTIFY_SECRET: "",
        NOTIFY_ENV: "production",
      },
    });
  },
});
