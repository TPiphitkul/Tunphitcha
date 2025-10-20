import json
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LogisticRegression
import joblib

# โหลด log
rows = [json.loads(line) for line in open("risk_log.jsonl")]
df = pd.DataFrame(rows)

# เตรียม feature
X = df[["req_per_min"]].copy()
X["path_login"] = df["path"].str.contains("/login").astype(int)
X["path_order"] = df["path"].str.contains("/order").astype(int)
y = df["attack"].astype(int)

# แบ่ง train/test
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# เทรน logistic regression
model = LogisticRegression()
model.fit(X_train, y_train)

acc = model.score(X_test, y_test)
print(f"Accuracy: {acc:.2f}")

joblib.dump(model, "risk_model.pkl")                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 
print("Saved model to risk_model.pkl")
