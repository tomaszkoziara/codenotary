import React, { useState } from "react";
import axios from "axios";

const apiPort = process.env.REACT_APP_API_PORT;
const apiAddr = `http://localhost:${apiPort}/api/v0`;

function App() {
  const [searchAccountName, setSearchAccountName] = useState("");
  const [accountingInfo, setAccountingInfo] = useState([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState("");

  const [newRecord, setNewRecord] = useState({
    accountNumber: "",
    accountName: "",
    iban: "",
    address: "",
    amount: "",
    type: "receiving",
  });

  const fetchAccountingInfo = () => {
    if (!searchAccountName && searchAccountName.trim() === "") {
      return;
    }
    setError(null);
    axios
      .get(
        `${apiAddr}/accountinginfo?page=${page}&pageSize=${pageSize}&accountName=${searchAccountName}`
      )
      .then((response) => {
        setAccountingInfo(response.data);
      })
      .catch((error) => {
        if (error.response && error.response.status === 404) {
          setAccountingInfo([]);
        } else {
          handleError(error);
        }
      });
  };

  const handleError = (error) => {
    let errorMessage = "An error occurred";

    if (error.response) {
      errorMessage =
        error.response.data?.message ||
        `Error: ${error.response.status} ${error.response.statusText}`;
    } else if (error.request) {
      errorMessage = "No response received. Please check your connection.";
    } else {
      errorMessage = error.message;
    }

    setError(errorMessage);
  };

  const handleSearch = () => {
    setPage(1);
    fetchAccountingInfo();
  };

  const handlePrevious = () => {
    if (page > 1) {
      setPage((prev) => prev - 1);
      fetchAccountingInfo();
    }
  };

  const handleNext = () => {
    setPage((prev) => prev + 1);
    fetchAccountingInfo();
  };

  const handleCreateRecord = () => {
    axios
      .post("${apiAddr}/accountinginfo", newRecord)
      .then(() => {
        setNewRecord({
          accountNumber: "",
          accountName: "",
          iban: "",
          address: "",
          amount: "",
          type: "receiving",
        });
        setSuccessMessage("Record created successfully!");
        setTimeout(() => setSuccessMessage(""), 3000);
      })
      .catch((error) => {
        handleError(error);
      });
  };

  return (
    <div className="container mx-auto p-6">
      {successMessage && (
        <div className="alert alert-success mb-4">{successMessage}</div>
      )}

      <div className="bg-base-100 p-4 rounded-lg shadow-md">
        <h2 className="text-xl font-bold mb-4">Create Accounting Info</h2>
        <div className="flex flex-col gap-4 items-center">
          <input
            type="text"
            placeholder="Account Number"
            value={newRecord.accountNumber}
            onChange={(e) =>
              setNewRecord({ ...newRecord, accountNumber: e.target.value })
            }
            className="input input-bordered w-full"
          />
          <input
            type="text"
            placeholder="Account Name"
            value={newRecord.accountName}
            onChange={(e) =>
              setNewRecord({ ...newRecord, accountName: e.target.value })
            }
            className="input input-bordered w-full"
          />
          <input
            type="text"
            placeholder="IBAN"
            value={newRecord.iban}
            onChange={(e) =>
              setNewRecord({ ...newRecord, iban: e.target.value })
            }
            className="input input-bordered w-full"
          />
          <input
            type="text"
            placeholder="Address"
            value={newRecord.address}
            onChange={(e) =>
              setNewRecord({ ...newRecord, address: e.target.value })
            }
            className="input input-bordered w-full"
          />
          <input
            type="number"
            placeholder="Amount"
            value={newRecord.amount}
            onChange={(e) =>
              setNewRecord({ ...newRecord, amount: parseFloat(e.target.value) })
            }
            className="input input-bordered w-full"
          />
          <select
            value={newRecord.type}
            onChange={(e) =>
              setNewRecord({ ...newRecord, type: e.target.value })
            }
            className="select select-bordered w-full"
          >
            <option value="receiving">Receiving</option>
            <option value="sending">Sending</option>
          </select>
          <button
            onClick={handleCreateRecord}
            className="btn btn-primary w-full"
          >
            Create Record
          </button>
        </div>
      </div>

      <div className="bg-base-100 p-4 rounded-lg shadow-md mt-8">
        <h2 className="text-xl font-bold mb-4">Search Accounting Info</h2>
        <div className="flex items-center gap-4">
          <input
            type="text"
            placeholder="Search by Account Name"
            value={searchAccountName}
            onChange={(e) => setSearchAccountName(e.target.value)}
            className="input input-bordered flex-grow"
          />
          <button onClick={handleSearch} className="btn btn-accent">
            Search
          </button>
        </div>

        <table className="table w-full mt-4">
          <thead>
            <tr>
              <th>Account Number</th>
              <th>Account Name</th>
              <th>IBAN</th>
              <th>Address</th>
              <th>Amount</th>
              <th>Type</th>
            </tr>
          </thead>
          <tbody>
            {accountingInfo.length > 0 ? (
              accountingInfo.map((info, index) => (
                <tr key={index}>
                  <td>{info.accountNumber}</td>
                  <td>{info.accountName}</td>
                  <td>{info.iban}</td>
                  <td>{info.address}</td>
                  <td>{info.amount}</td>
                  <td>{info.type}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan="6" className="text-center">
                  No records found
                </td>
              </tr>
            )}
          </tbody>
        </table>

        <div className="flex justify-between mt-4">
          <button
            className="btn btn-secondary"
            onClick={handlePrevious}
            disabled={page === 1}
          >
            Previous
          </button>
          <span>Page {page}</span>
          <button className="btn btn-secondary" onClick={handleNext}>
            Next
          </button>
        </div>
      </div>

      {error && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">Error</h3>
            <p className="py-4">{error}</p>
            <div className="modal-action">
              <button onClick={() => setError(null)} className="btn btn-error">
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
