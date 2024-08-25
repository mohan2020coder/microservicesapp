import { useState } from "react";

// Define types for the payload and response data
interface BrokerPayload {
  action?: string;
  auth?: {
    email: string;
    password: string;
  };
  log?: {
    name: string;
    data: string;
  };
  mail?: {
    from: string;
    to: string;
    subject: string;
    message: string;
  };
}

interface BrokerResponse {
  message: string;
  error?: boolean;
}

export default function Home() {
  const [output, setOutput] = useState<string>("Output shows here...");
  const [sent, setSent] = useState<string>("Nothing sent yet...");
  const [received, setReceived] = useState<string>("Nothing received yet...");

  const handleFetch = async (url: string, payload: BrokerPayload) => {
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      // Ensure the response is in text format before parsing
      const text = await response.text();
      console.log("Raw response:", text);

      const data: BrokerResponse = JSON.parse(text);
      setSent(JSON.stringify(payload, undefined, 4));
      setReceived(JSON.stringify(data, undefined, 4));

      if (data.error) {
        setOutput((prevOutput) => `${prevOutput}<br><strong>Error:</strong> ${data.message}`);
      } else {
        setOutput((prevOutput) => `${prevOutput}<br><strong>Response from service:</strong> ${data.message}`);
      }
    }catch (error) {
      // Type checking to ensure error has a message property
      if (error instanceof Error) {
          setOutput((prevOutput) => `${prevOutput}<br><br>Error: ${error.message}`);
      } else {
          // Handle cases where error is not an instance of Error
          setOutput((prevOutput) => `${prevOutput}<br><br>Unknown error occurred.`);
      }
  }
  };

  const handleBrokerClick = () => {
    const payload: BrokerPayload = {}; // Adjust payload as needed
    handleFetch("http://localhost:8080", payload);
  };

  const handleAuthBrokerClick = () => {
    const payload: BrokerPayload = {
      action: "auth",
      auth: {
        email: "admin@example.com",
        password: "verysecret",
      }
    };
    handleFetch("http://localhost:8080", payload);
  };

  const handleLogClick = () => {
    const payload: BrokerPayload = {
      action: "log",
      log: {
        name: "event",
        data: "Some kind of data",
      }
    };
    handleFetch("http://localhost:8080", payload);
  };

  const handleMailClick = () => {
    const payload: BrokerPayload = {
      action: "mail",
      mail: {
        from: "me@example.com",
        to: "you@there.com",
        subject: "Test email",
        message: "Hello world!",
      }
    };
    handleFetch("http://localhost:8080", payload);
  };

  const handleLogGClick = () => {
    const payload: BrokerPayload = {
      action: "log",
      log: {
        name: "event",
        data: "Some kind of gRPC data",
      }
    };
    handleFetch("http://localhost:8080/log-grpc", payload);
  };

  return (
    <div className="container mx-auto p-6">
      <div className="text-center">
        <h1 className="text-4xl font-bold mt-10 mb-5 text-gray-800">Test Microservices</h1>
        <hr className="mb-10 border-gray-300" />
        <button
          id="brokerBtn"
          className="px-6 py-3 bg-gray-800 text-white font-semibold rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50"
          onClick={handleBrokerClick}
        >
          Test Broker
        </button>
        <button
          id="authBrokerBtn"
          className="px-6 py-3 bg-gray-800 text-white font-semibold rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50 mt-2"
          onClick={handleAuthBrokerClick}
        >
          Test Auth
        </button>
        <button
          id="logBtn"
          className="px-6 py-3 bg-gray-800 text-white font-semibold rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50 mt-2"
          onClick={handleLogClick}
        >
          Test Log
        </button>
        <button
          id="mailBtn"
          className="px-6 py-3 bg-gray-800 text-white font-semibold rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50 mt-2"
          onClick={handleMailClick}
        >
          Test Mail
        </button>
        <button
          id="logGBtn"
          className="px-6 py-3 bg-gray-800 text-white font-semibold rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50 mt-2"
          onClick={handleLogGClick}
        >
          Test gRPC Log
        </button>

        <div
          id="output"
          className="mt-10 p-5 border border-gray-300 rounded-lg bg-gray-50 text-gray-700"
          style={{ outline: "1px solid silver", padding: "2em" }}
          dangerouslySetInnerHTML={{ __html: output }}
        />
      </div>
      <div className="mt-10 grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h4 className="text-2xl font-semibold text-gray-800 mb-4">Sent</h4>
          <div className="p-5 border border-gray-300 rounded-lg bg-gray-50 text-gray-700">
            <pre id="payload">{sent}</pre>
          </div>
        </div>
        <div>
          <h4 className="text-2xl font-semibold text-gray-800 mb-4">Received</h4>
          <div className="p-5 border border-gray-300 rounded-lg bg-gray-50 text-gray-700">
            <pre id="received">{received}</pre>
          </div>
        </div>
      </div>
    </div>
  );
}
