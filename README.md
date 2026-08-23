# Project ReSet

## Team Members
- Olasubomi Williams
- TOriola Joshua Oluwatunmise
- Imoleoluwa Emmanuella


---

## 🚀 Live Demo

*   **Live Application:** https://reset-lugo.onrender.com/
*   **Backend API:** [Link to your live backend API endpoint URL, if separate]
*   **Recorded Demo:** [Link to your recorded demo explaining how your solution works using Loom].


---

## 🎯 The Problem

*Many people in their day to day lives are met with unexpected bills that are gotten from subscription services that they have registered for but they have not or cannot unsubscribe from. These adds up costs that goes deep into their wallets and makes them at the mercy of predatory corporations*

> How might we help people save unnecessary subscription costs more effectively?

## ✨ Our Solution

*Our solution is to provide a way out of this trap - Meet ReSet - A automated Subscription cancellation software that gives users the option to directly opt-out of any unnecessary and unwanted subscription services.*


---

## 🛠️ Tech Stack

*List the major technologies, frameworks, and platforms you used to build your project.*

*   **Frontend:** HTMX
*   **Backend:** GOlang for the Backend, Python for the AI/ML
*   **Database:** PostgreSQL for Long Term Storage and Redis for in-memory caching 
*   **Deployment:** Render for Code deployment, Upstash for serverless 
*   **AI/APIs:** KIMI.AI for Research, Paystack API for subscription and cancellation,

---

## ⚙️ How to Set Up and Run Locally (Optional)

*Here are the steps to running the project running on a local machine.*


1.  Clone the repository:
    ```bash
    git clone https://github.com/12sub/ReSet/tree/confluence1
    ```
2.  Navigate to the project directory:
    ```bash
    cd /cmd/api
    ```
3.  Install dependencies:
    ```bash
    sudo apt install golang-go
    pip install -r requirements.txt
    go mod tidy
    ```
4.  Create a `.env.local` file and add the necessary environment variables:
    ```
    DATABASE_URL=...
    API_KEY=...
    ```
5.  To test for the subscription plans, you can either use flutterwave, or paystack:
    ```
    Go to paystack dashboard
    Register/Signup for paystack
    Go to Settings -> API Keys & Webhooks -> Create New API Key
    Go to Recurring -> Plans -> New Plans to create a New Plan
    ```
    *You can use this to test how the subscription cancellation software would work*

## ⚙️ Future Plans

*After this hackathon, we plan to improve on this project by*

1. Setting up Bank statement ingestation and Email receipt scanning functionality
2. Implement a smart cancellation engine
3. Setup renewal date tracking functionality
4. Setup a better AI Classification model
