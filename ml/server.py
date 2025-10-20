from flask import Flask, request, jsonify
import joblib
import numpy as np
import re

app = Flask(__name__)
model = joblib.load("risk_model.pkl")

@app.route("/predict", methods=["POST"])
def predict():
    data = request.json
    req_per_min = data.get("req_per_min", 0)
    path = data.get("path", "")
    path_login = 1 if re.search("/login", path) else 0
    path_order = 1 if re.search("/order", path) else 0
    X = np.array([[req_per_min, path_login, path_order]])
    prob = model.predict_proba(X)[0][1]
    risk_score = int(prob * 100)
    return jsonify({"risk_score": risk_score})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
