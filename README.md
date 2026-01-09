<div align="center">
  <img src="assets/logo.png" alt="DevTestrider Logo" width="128" height="128" />

  # DevTestrider

  **Next-Gen Real-Time Test Runner & Monitor for Go**

  [![Go Version](https://img.shields.io/github/go-mod/go-version/ismailtsdln/DevTestrider?style=flat-square)](https://golang.org)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
  [![Status](https://img.shields.io/badge/status-active-success.svg?style=flat-square)]()
</div>

---

**DevTestrider** is a powerful, developer-centric tool designed to supercharge your Go development workflow. It automatically watches your source code for changes, runs your tests in real-time, performs static analysis, and provides immediate feedback through a beautiful, modern Web UI and desktop notifications.

Stop switching context between your editor and terminal. Let DevTestrider handle the feedback loop.

## ✨ Key Features

*   **🚀 Real-Time Watcher**: Instantly detects file changes (recursive directory support) and triggers test runs.
*   **📊 Modern Web Dashboard**: A slick, responsive UI (React + Vite + TailwindCSS) displaying:
    *   Live test execution status.
    *   Detailed pass/fail breakdown per test case.
    *   Coverage trends and history.
    *   Static analysis issues.
*   **🛡️ Static Analysis Integration**: Automatically runs `go vet` to catch potential bugs and suspicious constructs alongside your tests.
*   **📄 Comprehensive Reporting**: Generates professional **HTML** and **PDF** reports for every run, perfect for archiving or sharing.
*   **🔔 Smart Notifications**: Native desktop notifications (MacOS/Linux/Windows) keep you informed without checking the UI.
*   **📈 Coverage Tracking**: Visual indicators for code coverage health (Green > 80%, Yellow > 50%, Red < 50%).
*   **🎨 CLI Experience**: Rich, color-coded terminal output using Lipgloss for those who prefer the command line.

## 🛠️ Installation

### Prerequisites
*   **Go**: 1.21 or higher
*   **Node.js**: 18+ (for building the frontend)

### Build form Source

```bash
# Clone the repository
git clone https://github.com/ismailtsdln/DevTestrider.git
cd DevTestrider

# Build the Frontend
cd web
npm install
npm run build
cd ..

# Build the Binary
go build -o devtestrider cmd/devtestrider/main.go
```

## 🚀 Usage

1.  **Initialize**: Create a configuration file (or use the default):
    ```yaml
    # testrider.yml
    watch:
      paths: ["."]
      ignore: [".git", "node_modules", "vendor"]
    report:
      formats: ["html", "pdf"]
      outputDir: "reports"
    notifications:
      enable: true
      channels: ["desktop"]
    server:
      port: 8085
    ```

2.  **Start**: Run the tool in your project root:
    ```bash
    ./devtestrider start
    ```

3.  **Monitor**: 
    *   Open your browser at `http://localhost:8085` to view the dashboard.
    *   Updates will stream in real-time as you code.

## 🧩 Architecture

DevTestrider is built with a modular architecture:

*   **Orchestrator**: Manages the lifecycle of the pipeline (Watcher -> Runner -> Analyzer -> Reporter -> Notifier).
*   **Engine**:
    *   **Watcher**: `fsnotify` based recursive file monitoring.
    *   **Runner**: Wraps `go test -json` for structured output.
    *   **Analyzer**: Wraps `go vet` for static analysis.
*   **Server**: Go HTTP server with Server-Sent Events (SSE) for real-time frontend updates.
*   **Report**: specialized engines for HTML (Text Templates) and PDF (Maroto) generation.
*   **Web**: Single Page Application built with React, TypeScript, and Recharts.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

<div align="center">
  <p>Built with ❤️ by Ismail Tasdelen</p>
</div>
