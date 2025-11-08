import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:rupu/presentation/views/horizontal_properties/controllers/live_voting_controller.dart';

class LiveVotingView extends GetWidget<LiveVotingController> {
  final int propertyId;
  final int votingGroupId;
  final int votingId;

  const LiveVotingView({
    super.key,
    required this.propertyId,
    required this.votingGroupId,
    required this.votingId,
  });

  @override
  String? get tag =>
      LiveVotingController.tagFor(propertyId, votingGroupId, votingId);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Votación en vivo'),
      ),
      body: Obx(
        () {
          if (controller.isLoading.value) {
            return const Center(child: CircularProgressIndicator());
          }
          if (controller.errorMessage.value != null) {
            return Center(child: Text(controller.errorMessage.value!));
          }
          return GridView.builder(
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
            ),
            itemCount: controller.units.length,
            itemBuilder: (context, index) {
              final unit = controller.units[index];
              return Card(
                color: unit.hasVoted ? Colors.green : Colors.white,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(unit.propertyUnitNumber),
                    Text(unit.residentName ?? 'Sin residente'),
                  ],
                ),
              );
            },
          );
        },
      ),
    );
  }
}
