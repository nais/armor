import React, {useState, useEffect} from 'react'
import {DataGrid} from '@mui/x-data-grid'

const columns = [
    {field: 'name', headerName: 'Policy Name', width: 200},
    {field: 'description', headerName: 'Description', width: 400},
    {field: 'fingerprint', headerName: 'Fingerprint', width: 150},
    {field: 'type', headerName: 'Policy Type', width: 150},
    {field: 'creation_timestamp', headerName: 'Creation', width: 250},
    {field: 'rules', headerName: 'Rules', width: 250},
]

const Grid = () => {

    const [tableData, setTableData] = useState([])

    useEffect(() => {
        fetch("http://localhost:8080/projects/plattformsikkerhet-dev-496e/policies")
            .then((data) => data.json())
            .then((data) => setTableData(data))

    }, [])

    return (
        <div style={{height: 900, width: '100%'}}>
            <DataGrid
                rows={tableData}
                columns={columns}
                checkboxSelection={true}
            />
        </div>
    )
}

export default Grid
