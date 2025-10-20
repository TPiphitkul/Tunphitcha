FROM python:3.10-slim
WORKDIR /app
COPY . .
RUN pip install --no-cache-dir flask joblib scikit-learn pandas numpy
EXPOSE 5000
CMD ["python", "server.py"]
