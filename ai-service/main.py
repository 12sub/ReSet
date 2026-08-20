"""
FastAPI service for subscription detection.

Run locally with:
    uvicorn main:app --reload --port 8001

Endpoints:
    POST /classify-transaction   -> classify a single transaction string
    POST /classify-batch         -> classify a list of transaction strings
    GET  /health                 -> simple health check
"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

from classifier import classify_transaction, classify_batch

app = FastAPI(title="ReSet Subscription Detection Service")

# Allow requests from any origin for the hackathon demo (Go backend,
# frontend, browser testing, etc). Fine for a demo; would be locked
# down to specific origins in production.
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=False,
    allow_methods=["*"],
    allow_headers=["*"],
)


class TransactionRequest(BaseModel):
    text: str


class BatchTransactionRequest(BaseModel):
    texts: list[str]


@app.get("/")
def root():
    return {"service": "ReSet Subscription Detection", "status": "running"}


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/classify-transaction")
def classify_single(req: TransactionRequest):
    return classify_transaction(req.text)


@app.post("/classify-batch")
def classify_many(req: BatchTransactionRequest):
    return {"results": classify_batch(req.texts)}