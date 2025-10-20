"""
Unit tests for transform_sensores_to_records function.
Tests the conversion from sensor-grouped format (backend) to flat timestamp format (ML model).
"""

import pytest
from typing import List, Dict, Any
from app.utils import transform_sensores_to_records


class TestTransformSensoresToRecords:
    """Test suite for sensor data format transformation"""
    
    def test_basic_transformation(self):
        """Test basic conversion with 2 sensors and 2 timestamps"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T15:59:59Z", "valor": 31.99},
                    {"tiempo": "2024-06-11T15:59:58Z", "valor": 31.98}
                ]
            },
            {
                "sensor_id": "A_Pres.PV",
                "datos": [
                    {"tiempo": "2024-06-11T15:59:59Z", "valor": 0.49},
                    {"tiempo": "2024-06-11T15:59:58Z", "valor": 0.48}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should have 2 records (one per timestamp)
        assert len(result) == 2
        
        # Check first record (sorted by timestamp, oldest first)
        assert result[0]["timestamp"] == "2024-06-11T15:59:58Z"
        assert result[0]["A_Temp.PV"] == 31.98
        assert result[0]["A_Pres.PV"] == 0.48
        
        # Check second record
        assert result[1]["timestamp"] == "2024-06-11T15:59:59Z"
        assert result[1]["A_Temp.PV"] == 31.99
        assert result[1]["A_Pres.PV"] == 0.49
    
    def test_multiple_sensors(self):
        """Test with multiple sensors (realistic scenario)"""
        sensores = [
            {
                "sensor_id": "A_ACR_Mot.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 0.001735}
                ]
            },
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99}
                ]
            },
            {
                "sensor_id": "A_Pres.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 0.49}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        assert len(result) == 1
        assert result[0]["timestamp"] == "2024-06-11T16:00:00Z"
        assert result[0]["A_ACR_Mot.PV"] == 0.001735
        assert result[0]["A_Temp.PV"] == 31.99
        assert result[0]["A_Pres.PV"] == 0.49
    
    def test_missing_data_points(self):
        """Test with sensors having different timestamps (sparse data)"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99},
                    {"tiempo": "2024-06-11T16:00:01Z", "valor": 32.00},
                    {"tiempo": "2024-06-11T16:00:02Z", "valor": 32.01}
                ]
            },
            {
                "sensor_id": "A_Pres.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 0.49},
                    # Missing 16:00:01
                    {"tiempo": "2024-06-11T16:00:02Z", "valor": 0.51}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should have 3 records, but middle one missing A_Pres.PV
        assert len(result) == 3
        
        # First record has both sensors
        assert "A_Temp.PV" in result[0]
        assert "A_Pres.PV" in result[0]
        
        # Middle record only has A_Temp.PV
        assert "A_Temp.PV" in result[1]
        assert "A_Pres.PV" not in result[1]
        
        # Last record has both sensors
        assert "A_Temp.PV" in result[2]
        assert "A_Pres.PV" in result[2]
    
    def test_timestamp_sorting(self):
        """Test that results are sorted by timestamp (oldest first)"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:02Z", "valor": 32.01},
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99},
                    {"tiempo": "2024-06-11T16:00:01Z", "valor": 32.00}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should be sorted in ascending order
        assert result[0]["timestamp"] == "2024-06-11T16:00:00Z"
        assert result[1]["timestamp"] == "2024-06-11T16:00:01Z"
        assert result[2]["timestamp"] == "2024-06-11T16:00:02Z"
    
    def test_empty_input(self):
        """Test with empty sensor list"""
        result = transform_sensores_to_records([])
        assert result == []
    
    def test_sensor_without_datos(self):
        """Test with sensor that has no data points"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": []
            },
            {
                "sensor_id": "A_Pres.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 0.49}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should only have record from A_Pres.PV
        assert len(result) == 1
        assert result[0]["A_Pres.PV"] == 0.49
        assert "A_Temp.PV" not in result[0]
    
    def test_missing_sensor_id(self):
        """Test handling of sensor without sensor_id"""
        sensores = [
            {
                # Missing sensor_id
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99}
                ]
            },
            {
                "sensor_id": "A_Pres.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 0.49}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should only process sensor with valid sensor_id
        assert len(result) == 1
        assert "A_Pres.PV" in result[0]
    
    def test_missing_timestamp(self):
        """Test handling of data point without timestamp"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99},
                    # Missing tiempo
                    {"valor": 32.00}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should only process record with valid timestamp
        assert len(result) == 1
        assert result[0]["timestamp"] == "2024-06-11T16:00:00Z"
    
    def test_missing_valor(self):
        """Test handling of data point without valor"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99},
                    # Missing valor
                    {"tiempo": "2024-06-11T16:00:01Z"}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should only process record with valid valor
        assert len(result) == 1
        assert result[0]["timestamp"] == "2024-06-11T16:00:00Z"
    
    def test_zero_valor(self):
        """Test that zero values are handled correctly (not treated as None)"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 0.0}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Zero is valid, should be included
        assert len(result) == 1
        assert result[0]["A_Temp.PV"] == 0.0
    
    def test_negative_valor(self):
        """Test that negative values are handled correctly"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": -5.5}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        assert len(result) == 1
        assert result[0]["A_Temp.PV"] == -5.5
    
    def test_duplicate_timestamps_same_sensor(self):
        """Test handling of duplicate timestamps in same sensor (should keep last)"""
        sensores = [
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99},
                    {"tiempo": "2024-06-11T16:00:00Z", "valor": 32.00}  # Duplicate
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should have 1 record with last value
        assert len(result) == 1
        assert result[0]["A_Temp.PV"] == 32.00  # Last value wins
    
    def test_real_world_scenario(self):
        """Test with real-world data structure from backend"""
        sensores = [
            {
                "sensor_id": "A_ACR_Mot.PV",
                "datos": [
                    {"tiempo": "2024-06-11T15:59:58Z", "valor": 0.001735421},
                    {"tiempo": "2024-06-11T15:59:59Z", "valor": 0.001740000}
                ]
            },
            {
                "sensor_id": "A_ACR_Mot.SV",
                "datos": [
                    {"tiempo": "2024-06-11T15:59:58Z", "valor": 0.001122334},
                    {"tiempo": "2024-06-11T15:59:59Z", "valor": 0.001125000}
                ]
            },
            {
                "sensor_id": "A_Temp.PV",
                "datos": [
                    {"tiempo": "2024-06-11T15:59:58Z", "valor": 31.98942017},
                    {"tiempo": "2024-06-11T15:59:59Z", "valor": 31.99942017}
                ]
            },
            {
                "sensor_id": "A_Pres.PV",
                "datos": [
                    {"tiempo": "2024-06-11T15:59:58Z", "valor": 0.488},
                    {"tiempo": "2024-06-11T15:59:59Z", "valor": 0.490}
                ]
            }
        ]
        
        result = transform_sensores_to_records(sensores)
        
        # Should have 2 records with all 4 sensors
        assert len(result) == 2
        
        # Check structure of first record
        first_record = result[0]
        assert first_record["timestamp"] == "2024-06-11T15:59:58Z"
        assert len(first_record) == 5  # timestamp + 4 sensors
        assert "A_ACR_Mot.PV" in first_record
        assert "A_ACR_Mot.SV" in first_record
        assert "A_Temp.PV" in first_record
        assert "A_Pres.PV" in first_record
        
        # Check structure of second record
        second_record = result[1]
        assert second_record["timestamp"] == "2024-06-11T15:59:59Z"
        assert len(second_record) == 5
        assert second_record["A_Temp.PV"] == 31.99942017


def test_integration_with_pumpwindow():
    """Test that transformed data can be used directly with PumpWindow model"""
    from app.main import PumpWindow
    
    sensores = [
        {
            "sensor_id": "A_Temp.PV",
            "datos": [
                {"tiempo": "2024-06-11T16:00:00Z", "valor": 31.99}
            ]
        }
    ]
    
    records = transform_sensores_to_records(sensores)
    
    # Should be able to create PumpWindow with transformed data
    pump_window = PumpWindow(pump="test_pump", records=records)
    
    assert pump_window.pump == "test_pump"
    assert len(pump_window.records) == 1
    assert pump_window.records[0]["timestamp"] == "2024-06-11T16:00:00Z"
    assert pump_window.records[0]["A_Temp.PV"] == 31.99


if __name__ == "__main__":
    # Run tests with pytest
    pytest.main([__file__, "-v", "--tb=short"])